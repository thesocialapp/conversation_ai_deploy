package communication

import (
	"bytes"
	"encoding/json"
	"net/http"
	"io" // Ensure to import "io"
	"go.uber.org/zap" // Ensure to import "go.uber.org/zap"
)

import "github.com/my/project/config"

var logger, _ = zap.NewProduction()

type Response struct {
	Response string `json:"response"`
}

func SendTextForProcessing(text string) (string, error) {
	url := config.Config.PythonAPIURL + "/process-text/"

	data := bytes.NewBuffer([]byte(text))
	resp, err := http.Post(url, "application/json", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseObj Response
	err = json.Unmarshal(body, &responseObj)
	if err != nil {
		logger.Error("Error unmarshalling JSON", zap.Error(err)) // Log the error
		return "", err
	}

	return responseObj.Response, nil
}

func ProcessText(text string) (string, error) {
	url := config.Config.PythonAPIURL + "/process"

	reqBody := map[string]string{"text": text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]string
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("Error unmarshalling JSON", zap.Error(err)) // Log the error
		return "", err
	}

	return result["response"], nil
}
