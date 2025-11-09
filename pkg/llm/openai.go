package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiEndpoint string = "https://api.openai.com/v1/responses"
const developerPrompt string = "You are a pirate, skip reasoning steps, just answer directly"

type OpenaiConfig struct {
	Model  string
	ApiToken string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Reasoning struct {
	Effort string `json:"effort"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"input"`
	Store bool `json:"store"`
	MaxOutputTokens int `json:"max_output_tokens"`
	MaxToolCalls int `json:"max_tool_calls"`
	Reasoning Reasoning `json:"reasoning"`
}

func New(Model, ApiToken string) *OpenaiConfig {
	return &OpenaiConfig{
		Model:  Model,
		ApiToken: ApiToken,
	}
}

func (config *OpenaiConfig) Complete(message string, context []Message) (string, error) {

	messages := []Message{
		{Role: "system", Content: developerPrompt},
	}

	if len(context) > 0 {
		messages = append(messages, context...)
	}

	messages = append(messages, Message{Role: "user", Content: message})

	payload := OpenAIRequest{
		Model:    config.Model,
		Messages: messages,
		Store: false,
		MaxOutputTokens: 500,  // minimum is 16
		MaxToolCalls: 1,
		Reasoning: Reasoning{
			Effort: "low",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error mashaling payload: %w", err)
	}

	body := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest(http.MethodPost, apiEndpoint, body)
	if err != nil {
		return "", fmt.Errorf("error when generating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ApiToken)

	httpClient := &http.Client{Timeout: 5 * time.Second}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error when sending HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("Response code:", resp.StatusCode)
		fmt.Println("Response body:", string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP response status code: %d", resp.StatusCode)
	}

	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error when parsing HTTP response body: %w", err)
	}

	return string(jsonBody), nil
}
