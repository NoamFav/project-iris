package config

import (
	"os"
	"path/filepath"
)

type WhisperModel int

const (
	WhisperTiny WhisperModel = iota
	WhisperTinyEn
	WhisperBase
	WhisperBaseEn
	WhisperSmall
	WhisperSmallEn
	WhisperMedium
	WhisperMediumEn
	WhisperLargeV1
	WhisperLargeV2
	WhisperLargeV3
	WhisperLargeV3Turbo
)

func (m WhisperModel) Path() string {
	base := filepath.Join(os.Getenv("HOME"), ".models", "whisper")

	switch m {
	case WhisperTiny:
		return filepath.Join(base, "ggml-tiny.bin")
	case WhisperTinyEn:
		return filepath.Join(base, "ggml-tiny.en.bin")
	case WhisperBase:
		return filepath.Join(base, "ggml-base.bin")
	case WhisperBaseEn:
		return filepath.Join(base, "ggml-base.en.bin")
	case WhisperSmall:
		return filepath.Join(base, "ggml-small.bin")
	case WhisperSmallEn:
		return filepath.Join(base, "ggml-small.en.bin")
	case WhisperMedium:
		return filepath.Join(base, "ggml-medium.bin")
	case WhisperMediumEn:
		return filepath.Join(base, "ggml-medium.en.bin")
	case WhisperLargeV1:
		return filepath.Join(base, "ggml-large-v1.bin")
	case WhisperLargeV2:
		return filepath.Join(base, "ggml-large-v2.bin")
	case WhisperLargeV3:
		return filepath.Join(base, "ggml-large-v3.bin")
	case WhisperLargeV3Turbo:
		return filepath.Join(base, "ggml-large-v3-turbo.bin")
	default:
		return ""
	}
}

func (m WhisperModel) String() string {
	switch m {
	case WhisperTiny:
		return "tiny"
	case WhisperTinyEn:
		return "tiny.en"
	case WhisperBase:
		return "base"
	case WhisperBaseEn:
		return "base.en"
	case WhisperSmall:
		return "small"
	case WhisperSmallEn:
		return "small.en"
	case WhisperMedium:
		return "medium"
	case WhisperMediumEn:
		return "medium.en"
	case WhisperLargeV1:
		return "large-v1"
	case WhisperLargeV2:
		return "large-v2"
	case WhisperLargeV3:
		return "large-v3"
	case WhisperLargeV3Turbo:
		return "large-v3-turbo"
	default:
		return "unknown"
	}
}
