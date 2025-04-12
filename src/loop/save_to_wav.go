package loop

import (
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"os"
)

func Float32ToInt16(samples []float32) []int {
	intSamples := make([]int, len(samples))
	for i, f := range samples {
		intSamples[i] = int(f * 32767)
	}
	return intSamples
}

func SaveAsWav(path string, floatSamples []float32, sampleRate int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := wav.NewEncoder(f, sampleRate, 16, 1, 1)

	intBuf := &audio.IntBuffer{
		Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1},
		Data:   Float32ToInt16(floatSamples),
	}

	if err := enc.Write(intBuf); err != nil {
		return err
	}
	return enc.Close()
}
