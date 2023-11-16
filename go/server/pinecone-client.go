package server

import (
	"net/http"
	"time"
)

const (
	baseURL = "https://api.openai.com/v1"
)

type PineConeClient struct {
	ApiKey     string
	HttpClient *http.Client
}

type PineClientOption func(*PineConeClient)

func WithTimeOut(timeout time.Duration) PineClientOption {
	return func(p *PineConeClient) {
		p.HttpClient.Timeout = 10
	}
}

func NewPineConeClient(apiKey string, options ...PineClientOption) *PineConeClient {
	client := &PineConeClient{
		ApiKey:     apiKey,
		HttpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}
