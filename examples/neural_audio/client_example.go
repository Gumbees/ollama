package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client_example <ollama_url>")
		fmt.Println("Example: client_example http://localhost:11434")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	
	fmt.Println("🎵 Neural Audio API Client Example 🎵")
	fmt.Println("=====================================")

	// Example 1: Basic Neural Audio with Generate API
	fmt.Println("\n1️⃣  Generate API with Neural Audio")
	runGenerateWithNeuralAudio(baseURL)

	// Example 2: Chat API with Neural Audio  
	fmt.Println("\n2️⃣  Chat API with Neural Audio")
	runChatWithNeuralAudio(baseURL)

	// Example 3: Advanced Configuration
	fmt.Println("\n3️⃣  Advanced Neural Audio Configuration")
	runAdvancedNeuralAudio(baseURL)
}

func runGenerateWithNeuralAudio(baseURL string) {
	// Basic neural audio request
	req := api.GenerateRequest{
		Model:  "llama3.2",
		Prompt: "Explain how neural networks work in simple terms.",
		Stream: boolPtr(true),
		
		// Enable neural audio with default settings
		NeuralAudio: boolPtr(true),
	}

	fmt.Printf("   📤 Sending request to %s/api/generate\n", baseURL)
	
	if err := streamRequest(baseURL+"/api/generate", req); err != nil {
		log.Printf("   ❌ Error: %v", err)
	}
}

func runChatWithNeuralAudio(baseURL string) {
	// Chat request with neural audio
	req := api.ChatRequest{
		Model: "llama3.2",
		Messages: []api.Message{
			{Role: "user", Content: "Tell me about quantum computing and how it differs from classical computing."},
		},
		Stream: boolPtr(true),
		
		// Enable neural audio with default settings
		NeuralAudio: boolPtr(true),
	}

	fmt.Printf("   📤 Sending request to %s/api/chat\n", baseURL)
	
	if err := streamChatRequest(baseURL+"/api/chat", req); err != nil {
		log.Printf("   ❌ Error: %v", err)
	}
}

func runAdvancedNeuralAudio(baseURL string) {
	// Advanced neural audio configuration
	req := api.GenerateRequest{
		Model:  "llama3.2",
		Prompt: "Write a creative story about an AI that creates music while thinking.",
		Stream: boolPtr(true),
		
		// Enable neural audio with custom configuration
		NeuralAudio: boolPtr(true),
		NeuralAudioConfig: &api.NeuralAudioConfig{
			LayerFilters: []string{
				"attention.query",
				"attention.key", 
				"attention.value",
				"feed_forward.up",
				"feed_forward.down",
			},
			MagnitudeToFreq:  true,
			VarianceToAmp:    true,
			SparsityToFilter: true,
			BaseFreq:         220.0, // Lower base frequency (A3)
			FreqRange:        [2]float64{110.0, 1760.0}, // Wider range (A2 to A6)
			AmpRange:         [2]float64{0.05, 0.95},
			UpdateRate:       "5ms", // Faster updates
			WaveType:         "sine",
			SampleRate:       44100,
		},
	}

	fmt.Printf("   📤 Sending advanced request to %s/api/generate\n", baseURL)
	
	if err := streamRequest(baseURL+"/api/generate", req); err != nil {
		log.Printf("   ❌ Error: %v", err)
	}
}

func streamRequest(url string, req api.GenerateRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("   📝 Streaming response with neural audio:")
	
	decoder := json.NewDecoder(resp.Body)
	for {
		var response api.GenerateResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}

		// Print text response
		if response.Response != "" {
			fmt.Printf("   💬 Text: %s", response.Response)
		}

		// Print neural audio information
		if response.NeuralAudio != nil {
			printNeuralAudioInfo(response.NeuralAudio)
		}

		if response.Done {
			fmt.Printf("\n   ✅ Complete! Duration: %v\n", response.TotalDuration)
			if response.NeuralAudio != nil {
				fmt.Printf("   🎵 Total audio samples generated: %d\n", len(response.NeuralAudio.AudioSamples))
			}
			break
		}
	}

	return nil
}

func streamChatRequest(url string, req api.ChatRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("   📝 Streaming chat response with neural audio:")
	
	decoder := json.NewDecoder(resp.Body)
	for {
		var response api.ChatResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}

		// Print message content
		if response.Message.Content != "" {
			fmt.Printf("   💬 %s: %s", response.Message.Role, response.Message.Content)
		}

		// Print neural audio information
		if response.NeuralAudio != nil {
			printNeuralAudioInfo(response.NeuralAudio)
		}

		if response.Done {
			fmt.Printf("\n   ✅ Complete! Duration: %v\n", response.TotalDuration)
			if response.NeuralAudio != nil {
				fmt.Printf("   🎵 Total audio samples generated: %d\n", len(response.NeuralAudio.AudioSamples))
			}
			break
		}
	}

	return nil
}

func printNeuralAudioInfo(audioData *api.NeuralAudioData) {
	if len(audioData.AudioSamples) > 0 {
		// Calculate some basic audio statistics
		var sum, min, max float32
		min = audioData.AudioSamples[0]
		max = audioData.AudioSamples[0]
		
		for _, sample := range audioData.AudioSamples {
			sum += sample
			if sample < min {
				min = sample
			}
			if sample > max {
				max = sample
			}
		}
		
		avg := sum / float32(len(audioData.AudioSamples))
		
		fmt.Printf("\n   🎵 Audio: %d samples, %.1fms, avg=%.3f, range=[%.3f,%.3f]", 
			len(audioData.AudioSamples), 
			audioData.Duration,
			avg, min, max)
	}

	// Print tensor statistics that contributed to this audio
	if len(audioData.TensorStats) > 0 {
		fmt.Printf("\n   🧠 Neural states:")
		for _, stat := range audioData.TensorStats {
			fmt.Printf("\n     • %s: mag=%.2f, var=%.2f, spar=%.2f → freq×%.2f, amp×%.2f",
				stat.Name,
				stat.Magnitude,
				stat.Variance, 
				stat.Sparsity,
				stat.FrequencyContribution,
				stat.AmplitudeContribution)
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}

/* Example Usage:

# Terminal 1: Start Ollama with neural audio support
ollama serve

# Terminal 2: Run the client example
go run client_example.go http://localhost:11434

Expected Output:
🎵 Neural Audio API Client Example 🎵
=====================================

1️⃣  Generate API with Neural Audio
   📤 Sending request to http://localhost:11434/api/generate
   📝 Streaming response with neural audio:
   💬 Text: Neural networks are computational models...
   🎵 Audio: 4410 samples, 100.0ms, avg=0.123, range=[-0.8,0.8]
   🧠 Neural states:
     • attention.query: mag=0.65, var=0.23, spar=0.15 → freq×0.70, amp×0.30
     • feed_forward.up: mag=0.42, var=0.31, spar=0.28 → freq×0.30, amp×0.70
   💬 Text: that are inspired by biological neural networks...
   🎵 Audio: 4410 samples, 100.0ms, avg=0.089, range=[-0.6,0.7]
   ...
   ✅ Complete! Duration: 2.3s
   🎵 Total audio samples generated: 88200

The client receives:
1. Regular text tokens as the AI generates them
2. Real-time audio data synchronized with the text generation
3. Tensor statistics showing which parts of the neural network contributed to each audio chunk
4. Metadata about audio properties (sample rate, duration, etc.)

This enables:
- Real-time "listening" to AI thinking
- Correlation between neural activity and generated text
- Audio visualization of model behavior
- Interactive AI experiences with audio feedback
*/
