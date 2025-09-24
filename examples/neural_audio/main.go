package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ollama/ollama/llama"
	"github.com/ollama/ollama/neural_audio"
)

func main() {
	// Set up logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	if len(os.Args) < 2 {
		fmt.Println("Usage: neural_audio <model_path>")
		fmt.Println("Example: neural_audio ./models/llama-2-7b-chat.q4_0.gguf")
		os.Exit(1)
	}

	modelPath := os.Args[1]
	
	fmt.Println("🎵 Neural Audio Synthesis for LLMs 🎵")
	fmt.Println("=====================================")
	fmt.Printf("Loading model: %s\n", modelPath)

	// Initialize llama backend
	llama.BackendInit()
	defer func() {
		fmt.Println("Cleaning up...")
	}()

	// Load the model
	model, err := loadModel(modelPath)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	defer llama.FreeModel(model)

	// Create context
	ctx, err := createContext(model)
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}

	fmt.Println("✅ Model loaded successfully")

	// Demo different neural audio configurations
	runNeuralAudioDemo(ctx)
}

func loadModel(modelPath string) (*llama.Model, error) {
	params := llama.ModelParams{
		NumGpuLayers: 32, // Adjust based on your GPU
		MainGpu:      0,
		UseMmap:      true,
		VocabOnly:    false,
		Progress: func(progress float32) {
			if int(progress*100)%10 == 0 {
				fmt.Printf("Loading progress: %.1f%%\n", progress*100)
			}
		},
	}

	return llama.LoadModelFromFile(modelPath, params)
}

func createContext(model *llama.Model) (*llama.Context, error) {
	contextParams := llama.NewContextParams(
		2048,  // context size
		512,   // batch size
		1,     // max sequences
		8,     // threads
		true,  // flash attention
		"f16", // kv cache type
	)

	return llama.NewContextWithModel(model, contextParams)
}

func runNeuralAudioDemo(ctx *llama.Context) {
	fmt.Println("\n🎼 Starting Neural Audio Demo")
	fmt.Println("============================")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Demo 1: Basic Neural Audio
	fmt.Println("\n1️⃣  Demo 1: Basic Neural Audio")
	fmt.Println("   Mapping tensor magnitude → frequency")
	fmt.Println("   Mapping tensor variance → amplitude")
	
	integration1, err := neural_audio.StartBasicNeuralAudio(ctx)
	if err != nil {
		log.Printf("Failed to start basic neural audio: %v", err)
	} else {
		// Run some inference to generate audio
		runInferenceDemo(ctx, "Hello, this is a test of neural audio synthesis. ")
		
		// Let it play for a bit
		time.Sleep(5 * time.Second)
		
		// Show stats
		stats := integration1.GetStats()
		fmt.Printf("   📊 Stats: %+v\n", stats)
		
		neural_audio.StopNeuralAudio(integration1)
	}

	// Demo 2: Advanced Neural Audio
	fmt.Println("\n2️⃣  Demo 2: Advanced Neural Audio")
	fmt.Println("   More complex mappings with layer filtering")
	fmt.Println("   Higher frequency range and faster updates")
	
	integration2, err := neural_audio.StartAdvancedNeuralAudio(ctx)
	if err != nil {
		log.Printf("Failed to start advanced neural audio: %v", err)
	} else {
		// Run more complex inference
		runInferenceDemo(ctx, "The quick brown fox jumps over the lazy dog. This sentence contains every letter of the alphabet. ")
		
		time.Sleep(5 * time.Second)
		
		stats := integration2.GetStats()
		fmt.Printf("   📊 Stats: %+v\n", stats)
		
		neural_audio.StopNeuralAudio(integration2)
	}

	// Demo 3: Custom Mapping
	fmt.Println("\n3️⃣  Demo 3: Custom Experimental Audio")
	fmt.Println("   Custom frequency ranges and waveforms")
	
	customMapping := neural_audio.TensorAudioMapping{
		LayerFilters: []string{
			"attention",
			"output",
		},
		MagnitudeToFreq:  true,
		VarianceToAmp:    true,
		SparsityToFilter: true,
		GradientToMod:    false,
		BaseFreq:         330.0, // E4 note
		FreqRange:        [2]float64{165.0, 1320.0}, // E3 to E6
		AmpRange:         [2]float64{0.2, 0.7},
		UpdateRate:       2 * time.Millisecond, // Very fast updates
	}
	
	integration3, err := neural_audio.CreateNeuralAudioSession(ctx, customMapping)
	if err != nil {
		log.Printf("Failed to start custom neural audio: %v", err)
	} else {
		// Run inference with a more complex prompt
		complexPrompt := `Explain the concept of neural networks and how they process information. 
Neural networks are computational models inspired by biological neural networks. 
They consist of interconnected nodes that process and transmit information.`
		
		runInferenceDemo(ctx, complexPrompt)
		
		time.Sleep(7 * time.Second)
		
		stats := integration3.GetStats()
		fmt.Printf("   📊 Stats: %+v\n", stats)
		
		neural_audio.StopNeuralAudio(integration3)
	}

	fmt.Println("\n🎵 Neural Audio Demo Complete!")
	fmt.Println("================================")
	fmt.Println("In a real implementation with actual audio output, you would have heard:")
	fmt.Println("• Different frequencies corresponding to neural activation patterns")
	fmt.Println("• Amplitude changes reflecting tensor variance") 
	fmt.Println("• Filter effects based on sparsity patterns")
	fmt.Println("• Real-time audio generation during model inference")
	
	// Wait for interrupt signal
	fmt.Println("\nPress Ctrl+C to exit...")
	<-sigChan
}

func runInferenceDemo(ctx *llama.Context, prompt string) {
	fmt.Printf("   🧠 Running inference: \"%s\"\n", prompt[:min(50, len(prompt))]+"...")
	
	// This is a simplified inference simulation
	// In reality, this would call the actual llama inference methods
	// which would trigger our neural audio callbacks
	
	// Simulate the inference process with some delay
	for i := 0; i < 20; i++ {
		// In real implementation, this would be:
		// - Tokenizing the prompt
		// - Running inference with ctx.Decode()
		// - Each decode call would trigger tensor operations
		// - Our callback would capture tensor states and generate audio
		
		time.Sleep(100 * time.Millisecond)
		if i%5 == 0 {
			fmt.Print(".")
		}
	}
	fmt.Println(" ✅")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Example of what the actual integration would look like in the main ollama codebase:

/*
In server/routes.go, the GenerateHandler would be modified like this:

func (s *Server) GenerateHandler(c *gin.Context) {
	// ... existing code ...
	
	// Check if neural audio is requested
	if req.NeuralAudio != nil && *req.NeuralAudio {
		mapping := neural_audio.DefaultAudioMapping()
		if req.NeuralAudioConfig != nil {
			// Use custom configuration
			mapping = *req.NeuralAudioConfig
		}
		
		audioIntegration, err := neural_audio.CreateNeuralAudioSession(r.Context(), mapping)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to enable neural audio: %v", err)})
			return
		}
		defer neural_audio.StopNeuralAudio(audioIntegration)
	}
	
	// ... rest of existing generation code ...
	// The neural audio would automatically capture tensor states during inference
}

And in api/types.go:

type GenerateRequest struct {
	// ... existing fields ...
	NeuralAudio       *bool                           `json:"neural_audio,omitempty"`
	NeuralAudioConfig *neural_audio.TensorAudioMapping `json:"neural_audio_config,omitempty"`
}
*/
