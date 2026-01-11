package openai

import (
	"net/http"
	"time"

	"github.com/gauthierdmn/shannon/pkg/llm"
)

type Client struct {
	model      string
	apiKey     string
	httpClient *http.Client
}

// create a new OpenAI API client
func NewClient(model, apiKey string) llm.Client {
	return &Client{
		model:      model,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}
