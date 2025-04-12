package loop

import (
	"fmt"
	"github.com/NoamFav/Iris/src/stt"
	"github.com/gordonklaus/portaudio"
	"log"
	"os/exec"
)

func StartLoop() {
	_ = exec.Command("rm", "-f", "audio/*.wav").Run()

	if err := portaudio.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer portaudio.Terminate()

	hostAPi, err := portaudio.DefaultHostApi()
	if err != nil {
		log.Fatal(err)
	}

	inputDevice := hostAPi.DefaultInputDevice
	fmt.Printf("\nUsing input device: %s\n", inputDevice.Name)

	const sampleRate = 44100
	const framesPerBuffer = 64
	buffer := make([]float32, framesPerBuffer)

	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   inputDevice,
			Channels: 1,
			Latency:  inputDevice.DefaultLowInputLatency,
		},
		SampleRate:      sampleRate,
		FramesPerBuffer: framesPerBuffer,
	}, buffer)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	fmt.Println(" Listening... Press Ctrl+C to stop.")

	collected := make([]float32, 0, sampleRate*5)
	chunkId := 0

	for {
		if err := stream.Read(); err != nil {
			if err == portaudio.InputOverflowed {
				fmt.Println(" Input overflow — skipping frame")
				continue
			}
			log.Fatal(err)
		}

		collected = append(collected, buffer...)

		if len(collected) >= sampleRate*5 {
			chunkFilename := fmt.Sprintf("audio/chunk_%03d.wav", chunkId)

			go func(data []float32, path string) {
				chunk := make([]float32, len(data))
				copy(chunk, data)

				err := SaveAsWav(path, chunk, sampleRate)
				if err != nil {
					fmt.Println(" Failed to save WAV:", err)
					return
				}

				text, err := stt.Transcribe(path)
				if err == nil {
					fmt.Println(" Whisper says:", text)
				} else {
					fmt.Println(" Transcription failed:", err)
				}
			}(collected, chunkFilename)

			collected = collected[:0]
			chunkId++
		}
	}
}
