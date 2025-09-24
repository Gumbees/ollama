// +build test

package neural_audio

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"testing"
	"time"
	"unsafe"
)

// TestAudioSynthesizer verifies the core audio synthesis functionality
func TestAudioSynthesizer(t *testing.T) {
	mapping := TensorAudioMapping{
		LayerFilters:     []string{"test_layer"},
		MagnitudeToFreq:  true,
		VarianceToAmp:    true,
		SparsityToFilter: false,
		BaseFreq:         440.0,
		FreqRange:        [2]float64{220.0, 880.0},
		AmpRange:         [2]float64{0.1, 0.8},
		UpdateRate:       10 * time.Millisecond,
	}

	synthesizer := NewNeuralAudioSynthesizer(mapping)
	if synthesizer == nil {
		t.Fatal("Failed to create synthesizer")
	}

	ctx := context.Background()
	err := synthesizer.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start synthesizer: %v", err)
	}

	// Simulate tensor processing
	for i := 0; i < 10; i++ {
		// Create a fake tensor pointer
		fakeTensor := unsafe.Pointer(uintptr(0x1000 + i))
		
		// First call - asking if we want to observe
		result := synthesizer.ProcessTensor(fakeTensor, true)
		if !result {
			t.Error("Expected ProcessTensor to return true for observation request")
		}

		// Second call - providing tensor data
		result = synthesizer.ProcessTensor(fakeTensor, false)
		if !result {
			t.Error("Expected ProcessTensor to return true for data provision")
		}
	}

	// Let it run for a bit
	time.Sleep(100 * time.Millisecond)

	stats := synthesizer.GetStats()
	if stats["total_tensors_processed"].(int64) < 10 {
		t.Errorf("Expected at least 10 tensors processed, got %v", stats["total_tensors_processed"])
	}

	synthesizer.Stop()
}

// TestTensorStatsExtraction verifies tensor statistics calculation
func TestTensorStatsExtraction(t *testing.T) {
	synthesizer := NewNeuralAudioSynthesizer(DefaultAudioMapping())
	
	// Test with a fake tensor
	fakeTensor := unsafe.Pointer(uintptr(0x2000))
	stats := synthesizer.extractTensorStats(fakeTensor)
	
	if stats.Name == "" {
		t.Error("Expected non-empty tensor name")
	}
	
	if stats.Magnitude < 0 || stats.Magnitude > 1 {
		t.Errorf("Expected magnitude in range [0,1], got %f", stats.Magnitude)
	}
	
	if stats.Variance < 0 || stats.Variance > 1 {
		t.Errorf("Expected variance in range [0,1], got %f", stats.Variance)
	}
	
	if stats.Sparsity < 0 || stats.Sparsity > 1 {
		t.Errorf("Expected sparsity in range [0,1], got %f", stats.Sparsity)
	}
}

// TestAudioParameterMapping verifies tensor stats are properly converted to audio
func TestAudioParameterMapping(t *testing.T) {
	mapping := TensorAudioMapping{
		MagnitudeToFreq:  true,
		VarianceToAmp:    true,
		SparsityToFilter: true,
		BaseFreq:         440.0,
		FreqRange:        [2]float64{220.0, 880.0},
		AmpRange:         [2]float64{0.1, 0.8},
		UpdateRate:       10 * time.Millisecond,
	}

	synthesizer := NewNeuralAudioSynthesizer(mapping)
	
	// Test with specific tensor stats
	stats := TensorStats{
		Name:      "test_tensor",
		Magnitude: 0.5, // Should map to middle of frequency range
		Variance:  0.3, // Should map to lower-middle amplitude
		Sparsity:  0.7, // High sparsity
	}
	
	// Capture the audio parameters (we'll need to modify updateAudioParams to return params for testing)
	originalChannelSize := len(synthesizer.audioChannel)
	synthesizer.updateAudioParams(stats)
	
	// Check that audio parameters were generated
	if len(synthesizer.audioChannel) <= originalChannelSize {
		t.Error("Expected audio parameters to be generated")
	}
}

// TestWaveformGeneration verifies different waveform types generate correctly
func TestWaveformGeneration(t *testing.T) {
	synthesizer := NewNeuralAudioSynthesizer(DefaultAudioMapping())
	
	testCases := []struct {
		waveType WaveType
		name     string
	}{
		{WaveSine, "sine"},
		{WaveSquare, "square"},
		{WaveSawtooth, "sawtooth"},
		{WaveTriangle, "triangle"},
		{WaveNoise, "noise"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := AudioParams{
				Frequency: 440.0,
				Amplitude: 0.5,
				Duration:  100, // 100ms
				WaveType:  tc.waveType,
			}
			
			synthesizer.generateAudioSample(params)
			
			// Check that audio was generated (buffer should have data)
			if len(synthesizer.audioBuffer) == 0 {
				t.Errorf("Expected audio buffer to contain data for %s wave", tc.name)
			}
			
			// Clear buffer for next test
			synthesizer.audioBuffer = synthesizer.audioBuffer[:0]
		})
	}
}

// TestMappingConfigurations verifies different mapping configurations work
func TestMappingConfigurations(t *testing.T) {
	configs := []struct {
		name    string
		mapping TensorAudioMapping
	}{
		{"default", DefaultAudioMapping()},
		{"advanced", AdvancedAudioMapping()},
		{"custom", TensorAudioMapping{
			MagnitudeToFreq:  false,
			VarianceToAmp:    true,
			SparsityToFilter: false,
			BaseFreq:         200.0,
			FreqRange:        [2]float64{100.0, 400.0},
			AmpRange:         [2]float64{0.05, 0.95},
			UpdateRate:       5 * time.Millisecond,
		}},
	}
	
	for _, config := range configs {
		t.Run(config.name, func(t *testing.T) {
			synthesizer := NewNeuralAudioSynthesizer(config.mapping)
			if synthesizer == nil {
				t.Fatalf("Failed to create synthesizer with %s mapping", config.name)
			}
			
			ctx := context.Background()
			err := synthesizer.Start(ctx)
			if err != nil {
				t.Fatalf("Failed to start synthesizer with %s mapping: %v", config.name, err)
			}
			
			synthesizer.Stop()
		})
	}
}

// BenchmarkTensorProcessing measures performance impact of tensor processing
func BenchmarkTensorProcessing(b *testing.B) {
	synthesizer := NewNeuralAudioSynthesizer(DefaultAudioMapping())
	ctx := context.Background()
	synthesizer.Start(ctx)
	defer synthesizer.Stop()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fakeTensor := unsafe.Pointer(uintptr(0x3000 + i))
		synthesizer.ProcessTensor(fakeTensor, true)
		synthesizer.ProcessTensor(fakeTensor, false)
	}
}

// BenchmarkAudioGeneration measures audio synthesis performance
func BenchmarkAudioGeneration(b *testing.B) {
	synthesizer := NewNeuralAudioSynthesizer(DefaultAudioMapping())
	
	params := AudioParams{
		Frequency: 440.0,
		Amplitude: 0.5,
		Duration:  50, // 50ms
		WaveType:  WaveSine,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		synthesizer.generateAudioSample(params)
	}
}

// TestConcurrentAccess verifies thread safety
func TestConcurrentAccess(t *testing.T) {
	synthesizer := NewNeuralAudioSynthesizer(DefaultAudioMapping())
	ctx := context.Background()
	err := synthesizer.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start synthesizer: %v", err)
	}
	defer synthesizer.Stop()
	
	// Start multiple goroutines processing tensors
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				fakeTensor := unsafe.Pointer(uintptr(0x4000 + id*1000 + j))
				synthesizer.ProcessTensor(fakeTensor, true)
				synthesizer.ProcessTensor(fakeTensor, false)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	stats := synthesizer.GetStats()
	totalProcessed := stats["total_tensors_processed"].(int64)
	if totalProcessed < 1000 {
		t.Errorf("Expected at least 1000 tensors processed, got %d", totalProcessed)
	}
}

// ExampleBasicUsage demonstrates basic usage
func ExampleBasicUsage() {
	// This example shows how neural audio would be used in practice
	// (Note: requires actual llama context for real usage)
	
	mapping := DefaultAudioMapping()
	synthesizer := NewNeuralAudioSynthesizer(mapping)
	
	ctx := context.Background()
	err := synthesizer.Start(ctx)
	if err != nil {
		slog.Error("Failed to start neural audio", "error", err)
		return
	}
	defer synthesizer.Stop()
	
	// Simulate some tensor processing
	for i := 0; i < 5; i++ {
		fakeTensor := unsafe.Pointer(uintptr(0x5000 + i))
		synthesizer.ProcessTensor(fakeTensor, true)
		synthesizer.ProcessTensor(fakeTensor, false)
		time.Sleep(10 * time.Millisecond)
	}
	
	stats := synthesizer.GetStats()
	fmt.Printf("Processed %d tensors, generated %d audio samples\n", 
		stats["total_tensors_processed"], 
		stats["audio_samples_generated"])
	
	// Output: Processed 5 tensors, generated 0 audio samples
}

// Helper function to validate audio parameters
func validateAudioParams(params AudioParams) error {
	if params.Frequency < 20 || params.Frequency > 20000 {
		return fmt.Errorf("frequency %f out of audible range", params.Frequency)
	}
	if params.Amplitude < 0 || params.Amplitude > 1 {
		return fmt.Errorf("amplitude %f out of valid range [0,1]", params.Amplitude)
	}
	if params.Duration < 0 {
		return fmt.Errorf("duration %d cannot be negative", params.Duration)
	}
	return nil
}

// Helper function to calculate audio quality metrics
func calculateAudioQuality(samples []float32) map[string]float64 {
	if len(samples) == 0 {
		return map[string]float64{"rms": 0, "peak": 0, "crest_factor": 0}
	}
	
	var sumSquares float64
	var peak float64
	
	for _, sample := range samples {
		abs := math.Abs(float64(sample))
		sumSquares += abs * abs
		if abs > peak {
			peak = abs
		}
	}
	
	rms := math.Sqrt(sumSquares / float64(len(samples)))
	crestFactor := peak / rms
	
	return map[string]float64{
		"rms":          rms,
		"peak":         peak,
		"crest_factor": crestFactor,
	}
}
