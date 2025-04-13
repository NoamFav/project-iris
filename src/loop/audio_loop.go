package loop

import (
	"fmt"
	"github.com/NoamFav/Iris/src/config"
	"github.com/NoamFav/Iris/src/stt"
	"github.com/gordonklaus/portaudio"
	"log"
	"os/exec"
	"time"
)

func StartLoop() {
	// Clean up any old audio chunks from previous runs
	if err := exec.Command("sh", "-c", "rm -f audio/*.wav").Run(); err != nil {
		fmt.Println("  Failed to clean audio dir:", err)
	}
	fmt.Println("  Cleaned up old chunks in audio/")

	// Initialize PortAudio (required before any audio operations)
	if err := portaudio.Initialize(); err != nil {
		log.Fatal("  Failed to initialize PortAudio:", err)
	}
	defer portaudio.Terminate() // Ensure cleanup when function exits
	fmt.Println("  PortAudio initialized")

	// Get default host API (needed to access default input device)
	hostAPi, err := portaudio.DefaultHostApi()
	if err != nil {
		log.Fatal("  Could not get default host API:", err)
	}

	// Use the system's default input device (like a microphone)
	inputDevice := hostAPi.DefaultInputDevice
	fmt.Printf("󱡬  Using input device: %s\n", inputDevice.Name)

	// Prepare a buffer that will hold a small frame of audio data
	buffer := make([]float32, config.FramesPerBuffer)

	// Open the input stream with given parameters
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   inputDevice,
			Channels: 1, // mono audio
			Latency:  inputDevice.DefaultLowInputLatency,
		},
		SampleRate:      config.SampleRate,
		FramesPerBuffer: config.FramesPerBuffer,
	}, buffer)
	if err != nil {
		log.Fatal("  Failed to open stream:", err)
	}
	defer stream.Close()

	// Start the audio stream
	if err := stream.Start(); err != nil {
		log.Fatal("  Failed to start stream:", err)
	}
	defer stream.Stop()

	fmt.Println("  Listening... Press Ctrl+C to stop.")

	// A rolling buffer to keep the last 10 seconds of audio
	ringBuffer := make([]float32, 0, config.SampleRate*10)

	// Timestamp for the last time we extracted a 5s chunk
	lastChunkTime := time.Now()

	// Chunk ID for naming output files
	chunkId := 0

	// Main audio loop
	for {
		// Read the next frame of audio into the buffer
		if err := stream.Read(); err != nil {
			if err == portaudio.InputOverflowed {
				fmt.Println("  Input overflow — skipping frame")
				continue
			}
			log.Fatal("  Stream read error:", err)
		}

		// Append new audio samples to the rolling buffer
		ringBuffer = append(ringBuffer, buffer...)

		// Trim the ring buffer to keep only the last 10s of audio
		if len(ringBuffer) > config.SampleRate*10 {
			ringBuffer = ringBuffer[len(ringBuffer)-config.SampleRate*10:]
		}

		// Every 5 seconds, process the latest 5 seconds of audio
		if time.Since(lastChunkTime) >= 5*time.Second && len(ringBuffer) >= config.SampleRate*5 {
			// Create a chunk of the last 5 seconds of audio
			chunk := make([]float32, config.SampleRate*5)
			copy(chunk, ringBuffer[len(ringBuffer)-config.SampleRate*5:])

			// Remember when the chunk began
			started := lastChunkTime
			lastChunkTime = time.Now()

			// Generate a file name for the chunk
			chunkFilename := fmt.Sprintf("%s/chunk_%03d.wav", config.AudioDir, chunkId)

			// Process the chunk concurrently to avoid blocking the loop
			go func(data []float32, path string, started time.Time) {
				// Save the audio chunk to WAV
				if err := SaveAsWav(path, data, config.SampleRate); err != nil {
					fmt.Println("  Failed to save WAV:", err)
					return
				}

				// Transcribe the WAV file using Whisper
				text, err := stt.Transcribe(path)
				if err == nil {
					fmt.Printf("  Whisper says: %q\n From: %s\n At: %s\n⇥ To: %s\n\n",
						text,
						path,
						started.Format("2006-01-02 15:04:05"),
						lastChunkTime.Format("2006-01-02 15:04:05"),
					)
				} else {
					fmt.Println("  Transcription failed:", err)
				}
			}(chunk, chunkFilename, started)

			// Increment chunk ID for the next audio chunk
			chunkId++
		}
	}
}
