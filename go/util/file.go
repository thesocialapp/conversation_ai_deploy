package util

import (
	"os"
)

func BytesToFile(data []byte, extension string) (string, error) {
	tempFile, err := os.CreateTemp("", "audio_*."+extension)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Write the audio bytes to the file
	_, err = tempFile.Write(data)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
