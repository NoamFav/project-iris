package stt

import (
	"bufio"
	"fmt"
	"github.com/NoamFav/Iris/src/config"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Transcribe(filePath string) (string, error) {

	modelPath := config.DefaultWhisperModel.Path()
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Fatalf("Whisper model not found: %s", modelPath)
	}

	cmd := exec.Command("whisper-cli", "-m", modelPath, "-f", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper error: %v\noutput: %s", err, string(out))
	}

	var transcription []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-->") {
			parts := strings.SplitN(line, "]", 2)
			if len(parts) == 2 {
				transcription = append(transcription, strings.TrimSpace(parts[1]))
			}
		}
	}

	return strings.Join(transcription, " "), nil
}
