/**
 * Neural Audio Chat Client
 * Handles real-time communication with Ollama API and neural audio processing
 */

class NeuralAudioClient {
    constructor(baseURL = 'http://localhost:11434') {
        this.baseURL = baseURL;
        this.audioContext = null;
        this.isConnected = false;
        this.currentStream = null;
    }

    async checkConnection() {
        try {
            const response = await fetch(`${this.baseURL}/api/version`);
            this.isConnected = response.ok;
            return this.isConnected;
        } catch (error) {
            this.isConnected = false;
            return false;
        }
    }

    async *streamChat(requestData) {
        try {
            const response = await fetch(`${this.baseURL}/api/chat`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder();

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);
                const lines = chunk.split('\n');

                for (const line of lines) {
                    if (line.trim()) {
                        try {
                            const data = JSON.parse(line);
                            yield data;
                        } catch (e) {
                            console.warn('Failed to parse JSON:', line, e);
                        }
                    }
                }
            }
        } catch (error) {
            console.error('Stream error:', error);
            throw error;
        }
    }

    async *streamGenerate(requestData) {
        try {
            const response = await fetch(`${this.baseURL}/api/generate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder();

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                const chunk = decoder.decode(value);
                const lines = chunk.split('\n');

                for (const line of lines) {
                    if (line.trim()) {
                        try {
                            const data = JSON.parse(line);
                            yield data;
                        } catch (e) {
                            console.warn('Failed to parse JSON:', line, e);
                        }
                    }
                }
            }
        } catch (error) {
            console.error('Stream error:', error);
            throw error;
        }
    }

    async initAudioContext() {
        if (!this.audioContext) {
            this.audioContext = new (window.AudioContext || window.webkitAudioContext)();
        }
        
        if (this.audioContext.state === 'suspended') {
            await this.audioContext.resume();
        }
        
        return this.audioContext;
    }

    async playAudioSamples(samples, sampleRate = 44100) {
        if (!samples || samples.length === 0) return;

        const audioContext = await this.initAudioContext();
        
        try {
            const audioBuffer = audioContext.createBuffer(1, samples.length, sampleRate);
            const channelData = audioBuffer.getChannelData(0);
            
            for (let i = 0; i < samples.length; i++) {
                channelData[i] = samples[i];
            }

            const source = audioContext.createBufferSource();
            source.buffer = audioBuffer;
            source.connect(audioContext.destination);
            source.start();
        } catch (error) {
            console.error('Audio playback error:', error);
        }
    }

    generateDefaultAudioConfig() {
        return {
            magnitude_to_freq: true,
            variance_to_amp: true,
            sparsity_to_filter: true,
            base_freq: 440.0,
            freq_range: [220.0, 880.0],
            amp_range: [0.1, 0.8],
            update_rate: "10ms",
            wave_type: "sine",
            sample_rate: 44100
        };
    }

    buildChatRequest(model, messages, neuralAudioEnabled, customConfig = {}) {
        const request = {
            model: model,
            messages: messages,
            stream: true
        };

        if (neuralAudioEnabled) {
            request.neural_audio = true;
            request.neural_audio_config = {
                ...this.generateDefaultAudioConfig(),
                ...customConfig
            };
        }

        return request;
    }

    buildGenerateRequest(model, prompt, neuralAudioEnabled, customConfig = {}) {
        const request = {
            model: model,
            prompt: prompt,
            stream: true
        };

        if (neuralAudioEnabled) {
            request.neural_audio = true;
            request.neural_audio_config = {
                ...this.generateDefaultAudioConfig(),
                ...customConfig
            };
        }

        return request;
    }
}

/**
 * Advanced Neural Audio Visualizer
 */
class NeuralAudioVisualizer {
    constructor(canvas) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d');
        this.animationId = null;
        this.audioData = [];
        this.frequencyData = [];
        this.tensorStats = [];
        this.isRunning = false;
        
        this.setupCanvas();
        this.setupColors();
    }

    setupCanvas() {
        const resizeCanvas = () => {
            const rect = this.canvas.getBoundingClientRect();
            this.canvas.width = rect.width * window.devicePixelRatio;
            this.canvas.height = rect.height * window.devicePixelRatio;
            this.ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
        };
        
        resizeCanvas();
        window.addEventListener('resize', resizeCanvas);
    }

    setupColors() {
        this.colors = {
            background: '#000',
            waveform: '#00ff88',
            frequency: '#3b82f6',
            attention: '#f59e0b',
            feedforward: '#ef4444',
            grid: '#333'
        };
    }

    start() {
        if (this.isRunning) return;
        this.isRunning = true;
        this.animate();
    }

    stop() {
        this.isRunning = false;
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }
    }

    animate() {
        if (!this.isRunning) return;
        
        this.draw();
        this.animationId = requestAnimationFrame(() => this.animate());
    }

    updateAudioData(samples) {
        this.audioData = samples || [];
        this.calculateFrequencyData();
    }

    updateTensorStats(stats) {
        this.tensorStats = stats || [];
    }

    calculateFrequencyData() {
        if (this.audioData.length === 0) {
            this.frequencyData = [];
            return;
        }

        // Simple FFT simulation for visualization
        const bands = 32;
        this.frequencyData = [];
        
        for (let i = 0; i < bands; i++) {
            // Simulate frequency content based on audio data
            let magnitude = 0;
            const samplesPerBand = Math.floor(this.audioData.length / bands);
            
            for (let j = 0; j < samplesPerBand; j++) {
                const index = i * samplesPerBand + j;
                if (index < this.audioData.length) {
                    magnitude += Math.abs(this.audioData[index]);
                }
            }
            
            magnitude = magnitude / samplesPerBand;
            this.frequencyData.push(magnitude);
        }
    }

    draw() {
        const width = this.canvas.width / window.devicePixelRatio;
        const height = this.canvas.height / window.devicePixelRatio;
        
        // Clear canvas
        this.ctx.fillStyle = this.colors.background;
        this.ctx.fillRect(0, 0, width, height);

        if (this.audioData.length > 0) {
            this.drawWaveform(width, height);
            this.drawFrequencySpectrum(width, height);
            this.drawTensorActivity(width, height);
        } else {
            this.drawIdleState(width, height);
        }

        this.drawGrid(width, height);
    }

    drawWaveform(width, height) {
        const waveformHeight = height * 0.3;
        const waveformY = height * 0.1;
        
        this.ctx.strokeStyle = this.colors.waveform;
        this.ctx.lineWidth = 2;
        this.ctx.beginPath();

        const sliceWidth = width / this.audioData.length;
        let x = 0;

        for (let i = 0; i < this.audioData.length; i++) {
            const v = this.audioData[i] * 0.5;
            const y = waveformY + (v + 1) * waveformHeight / 2;

            if (i === 0) {
                this.ctx.moveTo(x, y);
            } else {
                this.ctx.lineTo(x, y);
            }

            x += sliceWidth;
        }

        this.ctx.stroke();

        // Add label
        this.ctx.fillStyle = this.colors.waveform;
        this.ctx.font = '12px monospace';
        this.ctx.fillText('Waveform', 10, waveformY - 5);
    }

    drawFrequencySpectrum(width, height) {
        const spectrumHeight = height * 0.3;
        const spectrumY = height * 0.5;
        const barWidth = width / this.frequencyData.length;
        
        this.ctx.fillStyle = this.colors.frequency;
        
        for (let i = 0; i < this.frequencyData.length; i++) {
            const barHeight = this.frequencyData[i] * spectrumHeight;
            this.ctx.fillRect(
                i * barWidth, 
                spectrumY + spectrumHeight - barHeight, 
                barWidth - 2, 
                barHeight
            );
        }

        // Add label
        this.ctx.fillStyle = this.colors.frequency;
        this.ctx.font = '12px monospace';
        this.ctx.fillText('Frequency Spectrum', 10, spectrumY - 5);
    }

    drawTensorActivity(width, height) {
        const activityHeight = height * 0.15;
        const activityY = height * 0.85;
        
        if (this.tensorStats.length === 0) return;

        const barWidth = width / Math.max(this.tensorStats.length, 1);
        
        this.tensorStats.forEach((stat, index) => {
            const x = index * barWidth;
            
            // Draw magnitude bar
            this.ctx.fillStyle = this.colors.attention;
            const magHeight = stat.magnitude * activityHeight;
            this.ctx.fillRect(x, activityY + activityHeight - magHeight, barWidth * 0.3, magHeight);
            
            // Draw variance bar
            this.ctx.fillStyle = this.colors.feedforward;
            const varHeight = stat.variance * activityHeight;
            this.ctx.fillRect(x + barWidth * 0.35, activityY + activityHeight - varHeight, barWidth * 0.3, varHeight);
        });

        // Add labels
        this.ctx.fillStyle = '#fff';
        this.ctx.font = '10px monospace';
        this.ctx.fillText('Neural Activity', 10, activityY - 5);
        this.ctx.fillText('Mag', 15, activityY + activityHeight + 15);
        this.ctx.fillText('Var', 45, activityY + activityHeight + 15);
    }

    drawIdleState(width, height) {
        // Draw animated idle pattern
        const time = Date.now() * 0.001;
        const centerY = height / 2;
        const amplitude = height * 0.05;
        
        this.ctx.strokeStyle = this.colors.grid;
        this.ctx.lineWidth = 1;
        this.ctx.beginPath();
        
        for (let x = 0; x < width; x++) {
            const frequency1 = 0.02;
            const frequency2 = 0.05;
            const y = centerY + 
                Math.sin(x * frequency1 + time) * amplitude +
                Math.sin(x * frequency2 + time * 0.5) * amplitude * 0.5;
            
            if (x === 0) {
                this.ctx.moveTo(x, y);
            } else {
                this.ctx.lineTo(x, y);
            }
        }
        
        this.ctx.stroke();

        // Add idle text
        this.ctx.fillStyle = this.colors.grid;
        this.ctx.font = '16px monospace';
        this.ctx.textAlign = 'center';
        this.ctx.fillText('Waiting for neural activity...', width / 2, height / 2 - 30);
        this.ctx.textAlign = 'left';
    }

    drawGrid(width, height) {
        this.ctx.strokeStyle = this.colors.grid;
        this.ctx.lineWidth = 0.5;
        this.ctx.setLineDash([2, 2]);
        
        // Vertical lines
        const gridSpacing = width / 10;
        for (let x = gridSpacing; x < width; x += gridSpacing) {
            this.ctx.beginPath();
            this.ctx.moveTo(x, 0);
            this.ctx.lineTo(x, height);
            this.ctx.stroke();
        }
        
        // Horizontal lines
        const vGridSpacing = height / 6;
        for (let y = vGridSpacing; y < height; y += vGridSpacing) {
            this.ctx.beginPath();
            this.ctx.moveTo(0, y);
            this.ctx.lineTo(width, y);
            this.ctx.stroke();
        }
        
        this.ctx.setLineDash([]);
    }
}

/**
 * Audio Recording and Export Utilities
 */
class NeuralAudioRecorder {
    constructor() {
        this.isRecording = false;
        this.recordedChunks = [];
        this.startTime = null;
    }

    start() {
        this.isRecording = true;
        this.recordedChunks = [];
        this.startTime = Date.now();
        console.log('Started recording neural audio session');
    }

    stop() {
        this.isRecording = false;
        const duration = Date.now() - this.startTime;
        console.log(`Stopped recording. Duration: ${duration}ms, Chunks: ${this.recordedChunks.length}`);
        return this.recordedChunks;
    }

    addChunk(audioData) {
        if (!this.isRecording) return;
        
        this.recordedChunks.push({
            timestamp: Date.now() - this.startTime,
            audioSamples: audioData.audioSamples ? [...audioData.audioSamples] : [],
            tensorStats: audioData.tensorStats ? [...audioData.tensorStats] : [],
            sampleRate: audioData.sampleRate || 44100,
            duration: audioData.duration || 0
        });
    }

    exportSession() {
        if (this.recordedChunks.length === 0) {
            console.warn('No recorded data to export');
            return null;
        }

        const sessionData = {
            startTime: this.startTime,
            duration: Date.now() - this.startTime,
            chunks: this.recordedChunks,
            metadata: {
                totalSamples: this.recordedChunks.reduce((total, chunk) => total + chunk.audioSamples.length, 0),
                avgSampleRate: this.recordedChunks[0]?.sampleRate || 44100,
                tensorCount: this.recordedChunks.reduce((total, chunk) => total + chunk.tensorStats.length, 0)
            }
        };

        return sessionData;
    }

    downloadAsJSON() {
        const sessionData = this.exportSession();
        if (!sessionData) return;

        const blob = new Blob([JSON.stringify(sessionData, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = `neural-audio-session-${new Date().toISOString().slice(0, 19)}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        
        URL.revokeObjectURL(url);
    }

    downloadAsWAV() {
        // Combine all audio samples into a single WAV file
        const sessionData = this.exportSession();
        if (!sessionData) return;

        const allSamples = [];
        sessionData.chunks.forEach(chunk => {
            allSamples.push(...chunk.audioSamples);
        });

        const wavBlob = this.createWAVBlob(allSamples, sessionData.metadata.avgSampleRate);
        const url = URL.createObjectURL(wavBlob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = `neural-audio-${new Date().toISOString().slice(0, 19)}.wav`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        
        URL.revokeObjectURL(url);
    }

    createWAVBlob(samples, sampleRate) {
        const length = samples.length;
        const arrayBuffer = new ArrayBuffer(44 + length * 2);
        const view = new DataView(arrayBuffer);

        // WAV header
        const writeString = (offset, string) => {
            for (let i = 0; i < string.length; i++) {
                view.setUint8(offset + i, string.charCodeAt(i));
            }
        };

        writeString(0, 'RIFF');
        view.setUint32(4, 36 + length * 2, true);
        writeString(8, 'WAVE');
        writeString(12, 'fmt ');
        view.setUint32(16, 16, true);
        view.setUint16(20, 1, true);
        view.setUint16(22, 1, true);
        view.setUint32(24, sampleRate, true);
        view.setUint32(28, sampleRate * 2, true);
        view.setUint16(32, 2, true);
        view.setUint16(34, 16, true);
        writeString(36, 'data');
        view.setUint32(40, length * 2, true);

        // Convert samples to 16-bit PCM
        let offset = 44;
        for (let i = 0; i < length; i++) {
            const sample = Math.max(-1, Math.min(1, samples[i]));
            view.setInt16(offset, sample * 0x7FFF, true);
            offset += 2;
        }

        return new Blob([arrayBuffer], { type: 'audio/wav' });
    }
}

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { NeuralAudioClient, NeuralAudioVisualizer, NeuralAudioRecorder };
}
