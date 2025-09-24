package neural_audio

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"unsafe"

	"github.com/ollama/ollama/llama"
)

/*
#cgo CFLAGS: -std=c11
#cgo CXXFLAGS: -std=c++17
#cgo CPPFLAGS: -I${SRCDIR}/../llama/llama.cpp/include
#cgo CPPFLAGS: -I${SRCDIR}/../ml/backend/ggml/ggml/include

#include <stdlib.h>
#include "ggml.h"
#include "ggml-backend.h"

// Global storage for our callback instance
static void* g_neural_audio_instance = NULL;

// Forward declaration
bool neural_audio_eval_callback(struct ggml_tensor * tensor, bool ask, void * user_data);

// C wrapper to call our Go callback
bool neural_audio_eval_callback(struct ggml_tensor * tensor, bool ask, void * user_data) {
    // This calls back into Go code
    return neuralAudioEvalCallback(tensor, ask, user_data);
}

// Function to set the evaluation callback on a context
void set_neural_audio_callback(struct llama_context* ctx, void* instance) {
    g_neural_audio_instance = instance;
    // Note: This is conceptual - actual implementation would need to access
    // the internal scheduler from the llama context
    // ggml_backend_sched_set_eval_callback(sched, neural_audio_eval_callback, instance);
}
*/
import "C"

// Global registry to keep track of active synthesizers
var (
	synthesizerRegistry = make(map[unsafe.Pointer]*NeuralAudioSynthesizer)
	registryMutex       sync.RWMutex
)

// LlamaIntegration provides methods to integrate neural audio with llama contexts
type LlamaIntegration struct {
	synthesizer *NeuralAudioSynthesizer
	context     *llama.Context
	active      bool
}

// NewLlamaIntegration creates a new integration between a synthesizer and llama context
func NewLlamaIntegration(synthesizer *NeuralAudioSynthesizer, ctx *llama.Context) *LlamaIntegration {
	return &LlamaIntegration{
		synthesizer: synthesizer,
		context:     ctx,
	}
}

// Enable activates neural audio monitoring for this context
func (li *LlamaIntegration) Enable() error {
	if li.active {
		return fmt.Errorf("neural audio already enabled for this context")
	}

	// Register the synthesizer
	ctxPtr := unsafe.Pointer(li.context)
	registryMutex.Lock()
	synthesizerRegistry[ctxPtr] = li.synthesizer
	registryMutex.Unlock()

	// Set the callback on the context
	// Note: This is where we'd need to modify the llama context to expose
	// the internal scheduler and allow setting evaluation callbacks
	C.set_neural_audio_callback((*C.struct_llama_context)(ctxPtr), ctxPtr)

	li.active = true
	slog.Info("Neural audio enabled for llama context")
	return nil
}

// Disable deactivates neural audio monitoring
func (li *LlamaIntegration) Disable() {
	if !li.active {
		return
	}

	ctxPtr := unsafe.Pointer(li.context)
	registryMutex.Lock()
	delete(synthesizerRegistry, ctxPtr)
	registryMutex.Unlock()

	// Remove callback (would set to nil in real implementation)
	li.active = false
	slog.Info("Neural audio disabled for llama context")
}

// IsEnabled returns whether neural audio is currently active
func (li *LlamaIntegration) IsEnabled() bool {
	return li.active
}

//export neuralAudioEvalCallback
func neuralAudioEvalCallback(tensor *C.struct_ggml_tensor, ask C.bool, userData unsafe.Pointer) C.bool {
	registryMutex.RLock()
	synthesizer, exists := synthesizerRegistry[userData]
	registryMutex.RUnlock()

	if !exists {
		// No synthesizer registered for this context
		return C.bool(true)
	}

	// Call the synthesizer's process method
	result := synthesizer.ProcessTensor(unsafe.Pointer(tensor), bool(ask))
	return C.bool(result)
}

// DefaultAudioMapping provides a reasonable default mapping configuration
func DefaultAudioMapping() TensorAudioMapping {
	return TensorAudioMapping{
		LayerFilters: []string{
			"attention",
			"feed_forward", 
			"norm",
			"output",
		},
		MagnitudeToFreq:  true,
		VarianceToAmp:    true,
		SparsityToFilter: true,
		GradientToMod:    false, // Gradients not available during inference
		BaseFreq:         440.0, // A4 note
		FreqRange:        [2]float64{220.0, 880.0}, // A3 to A5
		AmpRange:         [2]float64{0.1, 0.8},
		UpdateRate:       10 * 1000000, // 10ms in nanoseconds
	}
}

// AdvancedAudioMapping provides a more complex mapping for experimentation
func AdvancedAudioMapping() TensorAudioMapping {
	return TensorAudioMapping{
		LayerFilters: []string{
			"attention.query",
			"attention.key", 
			"attention.value",
			"attention.output",
			"mlp.gate_proj",
			"mlp.up_proj",
			"mlp.down_proj",
		},
		MagnitudeToFreq:  true,
		VarianceToAmp:    true,
		SparsityToFilter: true,
		GradientToMod:    false,
		BaseFreq:         220.0,
		FreqRange:        [2]float64{110.0, 1760.0}, // A2 to A6 (wider range)
		AmpRange:         [2]float64{0.05, 0.9},
		UpdateRate:       5 * 1000000, // 5ms for more responsive audio
	}
}

// CreateNeuralAudioSession sets up a complete neural audio session
func CreateNeuralAudioSession(ctx *llama.Context, mapping TensorAudioMapping) (*LlamaIntegration, error) {
	// Create synthesizer
	synthesizer := NewNeuralAudioSynthesizer(mapping)
	
	// Create integration
	integration := NewLlamaIntegration(synthesizer, ctx)
	
	// Start the synthesizer
	audioCtx := context.Background()
	if err := synthesizer.Start(audioCtx); err != nil {
		return nil, fmt.Errorf("failed to start synthesizer: %w", err)
	}
	
	// Enable integration
	if err := integration.Enable(); err != nil {
		synthesizer.Stop()
		return nil, fmt.Errorf("failed to enable integration: %w", err)
	}
	
	slog.Info("Neural audio session created successfully")
	return integration, nil
}

// Example usage functions

// StartBasicNeuralAudio is a simple one-line function to enable neural audio
func StartBasicNeuralAudio(ctx *llama.Context) (*LlamaIntegration, error) {
	mapping := DefaultAudioMapping()
	return CreateNeuralAudioSession(ctx, mapping)
}

// StartAdvancedNeuralAudio enables neural audio with more detailed mappings
func StartAdvancedNeuralAudio(ctx *llama.Context) (*LlamaIntegration, error) {
	mapping := AdvancedAudioMapping()
	return CreateNeuralAudioSession(ctx, mapping)
}

// StopNeuralAudio cleanly shuts down a neural audio session
func StopNeuralAudio(integration *LlamaIntegration) {
	integration.Disable()
	if integration.synthesizer != nil {
		integration.synthesizer.Stop()
	}
	slog.Info("Neural audio session stopped")
}
