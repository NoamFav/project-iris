package loop

import (
	"fmt"
	"github.com/NoamFav/Iris/src/config"
	"github.com/NoamFav/Iris/src/stt"
	"github.com/gordonklaus/portaudio"
	"log"
	"os/exec"
)

func StartLoop() {
	_ = exec.Command("rm", "-f", "audio/*.wav").Run()
	fmt.Println("  Cleaned up old chunks in audio/")

	if err := portaudio.Initialize(); err != nil {
		log.Fatal("  Failed to initialize PortAudio:", err)
	}
	defer portaudio.Terminate()
	fmt.Println("  PortAudio initialized")

	hostAPi, err := portaudio.DefaultHostApi()
	if err != nil {
		log.Fatal("  Could not get default host API:", err)
	}

	inputDevice := hostAPi.DefaultInputDevice
	fmt.Printf("󱡬  Using input device: %s\n", inputDevice.Name)

	buffer := make([]float32, config.FramesPerBuffer)

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   inputDevice,
			Channels: 1,
			Latency:  inputDevice.DefaultLowInputLatency,
		},
		SampleRate:      config.SampleRate,
		FramesPerBuffer: config.FramesPerBuffer,
	}, buffer)
	if err != nil {
		log.Fatal("  Failed to open stream:", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatal("  Failed to start stream:", err)
	}
	defer stream.Stop()

	fmt.Println("  Listening... Press Ctrl+C to stop.")

	collected := make([]float32, 0, config.SampleRate*5)
	chunkId := 0

	for {
		if err := stream.Read(); err != nil {
			if err == portaudio.InputOverflowed {
				fmt.Println("  Input overflow — skipping frame")
				continue
			}
			log.Fatal("  Stream read error:", err)
		}

		collected = append(collected, buffer...)

		if len(collected) >= config.SampleRate*5 {
			chunkFilename := fmt.Sprintf("%s/chunk_%03d.wav", config.AudioDir, chunkId)
			fmt.Printf("  Chunk captured: %s\n", chunkFilename)

			go func(data []float32, path string) {
				chunk := make([]float32, len(data))
				copy(chunk, data)

				err := SaveAsWav(path, chunk, config.SampleRate)
				if err != nil {
					fmt.Println("  Failed to save WAV:", err)
					return
				}
				fmt.Println("  Saved WAV:", path)

				text, err := stt.Transcribe(path)
				if err == nil {
					fmt.Println("  Whisper says:", text)
				} else {
					fmt.Println("  Transcription failed:", err)
				}
			}(collected, chunkFilename)

			collected = collected[:0]
			chunkId++
		}
	}
}
