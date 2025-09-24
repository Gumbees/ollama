# 🎵 Ollama Neural Audio Web Interface

**A beautiful, real-time web interface for chatting with AI models while listening to their neural network activity!**

This web interface provides a complete GUI for Ollama's neural audio feature, allowing you to:
- Chat with AI models in real-time
- Listen to neural network "thoughts" as audio
- Visualize brain activity with live waveforms and spectrograms
- Record and export neural audio sessions
- Configure audio synthesis parameters

## 🚀 Quick Start

### 1. Start Ollama Server
```bash
# Start Ollama with the neural audio web interface enabled
ollama serve
```

### 2. Open Web Interface
Navigate to: **http://localhost:11434/**

You'll automatically be redirected to the neural audio interface at `/neural-audio`

### 3. Start Chatting with Neural Audio! 🎵
1. ✅ Ensure "🎵 Enable Neural Audio" is checked
2. 💬 Type a message and press Send
3. 🎧 Listen to the AI's neural network processing in real-time
4. 📊 Watch the live visualizations of neural activity

## 🎛️ Features

### **Real-Time Neural Audio Chat**
- **Live streaming chat** with any Ollama model
- **Neural audio synthesis** plays as the AI generates responses
- **Synchronized playback** - audio matches the exact neural states creating each word
- **Multi-model support** - works with llama3.2, qwen2.5, mistral, codellama, etc.

### **Advanced Audio Visualization**
- **Real-time waveform** display of neural activity
- **Frequency spectrum** analyzer showing neural patterns
- **Tensor activity bars** for individual neural network layers
- **Animated visualizations** that respond to AI thinking patterns

### **Configurable Audio Synthesis**
```javascript
// Customize how neural states become audio
{
  base_freq: 440,           // Base frequency (Hz)
  freq_range: [220, 880],   // Frequency range
  wave_type: "sine",        // Waveform: sine, square, sawtooth, triangle
  sample_rate: 44100,       // Audio quality
  update_rate: "10ms"       // How often audio updates
}
```

### **Audio Recording & Export**
- **🎙️ Record sessions** - capture entire neural audio conversations
- **📁 Export as JSON** - detailed neural state data for analysis
- **🎵 Export as WAV** - audio files you can play anywhere
- **📊 Session statistics** - tensor counts, audio samples, duration

### **Responsive Design**
- **💻 Desktop optimized** with dual-panel layout
- **📱 Mobile friendly** with collapsible panels
- **🎨 Beautiful gradients** and smooth animations
- **🌙 Dark theme** optimized for audio visualization

## 🎵 How Neural Audio Works

### **What You Hear**
Different aspects of the AI's "thinking" create different sounds:

1. **🧠 Attention Mechanisms**
   - **Query/Key/Value operations** → Distinct frequency patterns
   - **Multi-head attention** → Layered harmonic content
   - **Attention weights** → Volume changes

2. **⚡ Neural Network Layers**
   - **Feed-forward networks** → Sharp, rhythmic patterns
   - **Normalization layers** → Smooth, flowing tones
   - **Deeper layers** → Lower frequency ranges

3. **🤔 AI "Emotions"**
   - **Confidence** → Steady, clear tones
   - **Uncertainty** → Fluctuating, chaotic audio
   - **Complex reasoning** → Rich, complex harmonics
   - **Simple responses** → Clean, simple waveforms

### **Audio Mapping Examples**

```javascript
// Default Configuration - Balanced Audio
{
  magnitude_to_freq: true,    // Larger activations = higher pitch
  variance_to_amp: true,      // More chaos = louder volume
  sparsity_to_filter: true,   // Sparse patterns = filter effects
  base_freq: 440,             // A4 musical note
  freq_range: [220, 880]      // A3 to A5 range
}

// Experimental Configuration - Wide Range
{
  base_freq: 110,             // Lower base (A2)
  freq_range: [55, 1760],     // A1 to A6 (3 octaves!)
  wave_type: "sawtooth",      // Richer harmonics
  update_rate: "5ms"          // Super responsive
}
```

## 🛠️ Advanced Usage

### **Direct API Integration**
The web interface calls the standard Ollama API with neural audio enabled:

```javascript
// Chat API with Neural Audio
const response = await fetch('/api/chat', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    model: 'llama3.2',
    messages: [{ role: 'user', content: 'Tell me about AI' }],
    stream: true,
    neural_audio: true,
    neural_audio_config: {
      magnitude_to_freq: true,
      variance_to_amp: true,
      base_freq: 440
    }
  })
});

// Response includes both text and audio
// {"message":{"content":"AI is..."},"neural_audio":{"audio_samples":[0.1,-0.2,...],"tensor_stats":[...]}}
```

### **Custom Audio Processing**
```javascript
// Access the neural audio client directly
const client = new NeuralAudioClient('http://localhost:11434');

// Stream with custom processing
for await (const chunk of client.streamChat(request)) {
  if (chunk.neural_audio) {
    // Custom audio processing
    processNeuralAudio(chunk.neural_audio);
    
    // Play with custom effects
    await client.playAudioSamples(
      chunk.neural_audio.audio_samples,
      chunk.neural_audio.sample_rate
    );
  }
}
```

### **Recording Neural Sessions**
```javascript
const recorder = new NeuralAudioRecorder();

// Start recording
recorder.start();

// ... chat with AI (audio automatically recorded) ...

// Export session data
recorder.downloadAsJSON();  // Detailed analysis data
recorder.downloadAsWAV();   // Audio file
```

## 🎛️ Interface Controls

### **Chat Panel**
- **Model Selection** - Choose from available Ollama models
- **Neural Audio Toggle** - Enable/disable neural audio synthesis  
- **Audio Configuration** - Fine-tune synthesis parameters
- **Message History** - Full conversation with timestamps
- **Status Indicator** - Connection and generation status

### **Audio Panel**
- **Live Visualizer** - Real-time waveform and spectrum display
- **Audio Controls** - Mute, record, pause visualization
- **Neural Statistics** - Live tensor processing metrics
- **Tensor Activity** - Individual layer activity display

### **Keyboard Shortcuts**
- **Enter** - Send message
- **Shift+Enter** - New line in message
- **Ctrl+M** - Toggle mute
- **Ctrl+R** - Start/stop recording

## 🔧 Configuration Options

### **URL Parameters**
Customize the interface with URL parameters:

```
http://localhost:11434/neural-audio?
  model=llama3.2&
  neural_audio=true&
  base_freq=330&
  wave_type=square&
  auto_start=true
```

### **Environment Variables**
```bash
# Customize default settings
export OLLAMA_NEURAL_AUDIO_ENABLED=true
export OLLAMA_NEURAL_AUDIO_BASE_FREQ=440
export OLLAMA_NEURAL_AUDIO_SAMPLE_RATE=44100
```

### **Browser Storage**
Settings are automatically saved to localStorage:
- Audio configuration preferences
- Model selection
- Volume and mute state
- Visualization preferences

## 🎨 Customization

### **Themes**
Modify the CSS variables in `index.html`:

```css
:root {
  --neural-primary: #3b82f6;      /* Primary blue */
  --neural-secondary: #10b981;    /* Green accents */
  --neural-background: #667eea;   /* Background gradient */
  --neural-waveform: #00ff88;     /* Waveform color */
}
```

### **Audio Visualization**
Customize the visualizer in `neural-audio.js`:

```javascript
// Custom color schemes
this.colors = {
  background: '#000',
  waveform: '#00ff88',     // Matrix green
  frequency: '#3b82f6',    // Electric blue  
  attention: '#f59e0b',    // Attention orange
  feedforward: '#ef4444'   // Processing red
};
```

## 🚀 Performance Tips

### **Optimal Settings**
- **Sample Rate**: 44100 Hz for best quality, 22050 Hz for performance
- **Update Rate**: 10ms for smooth audio, 5ms for maximum responsiveness
- **Buffer Size**: Automatically managed for minimal latency

### **Browser Compatibility**
- **✅ Chrome/Edge** - Full Web Audio API support
- **✅ Firefox** - Full support with minor audio latency
- **✅ Safari** - Requires user interaction to start audio
- **⚠️ Mobile** - Limited by device audio capabilities

### **Network Optimization**
- **Local Ollama** - Best performance, <10ms audio latency
- **Remote Ollama** - Add network delay, still very usable
- **Streaming** - Chunked transfer for real-time experience

## 🐛 Troubleshooting

### **No Audio Playing**
1. Check browser audio permissions
2. Ensure Ollama server is running
3. Verify `neural_audio: true` in requests
4. Check browser console for errors

### **Visualization Not Working**
1. Verify Canvas API support
2. Check `neural-audio.js` is loaded
3. Ensure WebGL is enabled (optional)

### **Connection Issues**
1. Verify Ollama is running: `ollama list`
2. Check firewall settings
3. Try different browser/incognito mode
4. Check network connectivity

### **Performance Issues**
1. Lower sample rate to 22050 Hz
2. Increase update rate to 20ms
3. Disable visualization during recording
4. Close other audio applications

## 🎵 Example Sessions

### **"Explain Quantum Physics"**
- **🎵 Audio**: Complex, layered harmonics with rapid frequency changes
- **📊 Patterns**: High variance during complex explanations
- **🧠 Neural**: Heavy attention layer activity, moderate feed-forward

### **"Write a Simple Poem"**
- **🎵 Audio**: Rhythmic, musical patterns with steady beat
- **📊 Patterns**: Regular frequency changes matching verse structure  
- **🧠 Neural**: Balanced attention/feed-forward, creative spikes

### **"What's 2+2?"**
- **🎵 Audio**: Quick, simple tone with minimal variation
- **📊 Patterns**: Low variance, stable frequency
- **🧠 Neural**: Minimal processing, fast convergence

## 📚 API Reference

See the full neural audio API documentation in the main README:
`/neural_audio/README.md`

## 🤝 Contributing

The web interface is built with vanilla HTML/CSS/JavaScript for maximum compatibility:

- **index.html** - Main interface and styling
- **neural-audio.js** - Audio client and visualization classes
- **README.md** - This documentation

To contribute:
1. Test with different models and configurations
2. Add new visualization modes
3. Improve mobile responsiveness  
4. Add audio effects and filters
5. Create new themes and customizations

---

**🎵 Experience AI like never before - Listen to the neural symphony of artificial intelligence! 🧠**
