package stt

import (
	"bufio"
	"fmt"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func Transcribe(filePath string) (string, error) {

	usr, _ := user.Current()
	modelPath := filepath.Join(usr.HomeDir, ".models/voice/ggml-large-v3-turbo.bin")

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
