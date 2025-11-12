package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const apiEndpoint string = "https://api.openai.com/v1/responses"
const developerPrompt string = "You are a pirate, skip reasoning steps, just answer directly"

type OpenaiClient struct {
	model      string
	apiToken   string
	httpClient *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Reasoning struct {
	Effort string `json:"effort"`
}

type OpenaiRequest struct {
	Model           string    `json:"model"`
	Messages        []Message `json:"input"`
	Store           bool      `json:"store"`
	MaxOutputTokens int       `json:"max_output_tokens"`
	MaxToolCalls    int       `json:"max_tool_calls"`
	Reasoning       Reasoning `json:"reasoning"`
}

type OpenaiResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OpenaiResponseOutput struct {
	Type    string                        `json:"type"`
	Status  string                        `json:"status"`
	Content []OpenaiResponseOutputContent `json:"content"`
}

type OpenaiResponseOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type OpenaiResponse struct {
	Status string                 `json:"status"`
	Error  *OpenaiResponseError   `json:"error"`
	Output []OpenaiResponseOutput `json:"output"`
}

// create a new OpenAI client
func New(model, apiToken string) *OpenaiClient {
	return &OpenaiClient{
		model:      model,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// complete a text message using the Open AI completion API
func (client *OpenaiClient) Complete(ctx context.Context, message string, history []Message) (string, error) {

	req, err := client.buildRequest(ctx, message, history)
	if err != nil {
		return "", fmt.Errorf("error when building openai request: %s", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error when sending HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return parseResponse(resp)
}

// build an HTTP request for the OpenAI completion API
func (client *OpenaiClient) buildRequest(ctx context.Context, message string, history []Message) (*http.Request, error) {

	messages := appendHistory(message, history)

	payload := OpenaiRequest{
		Model:           client.model,
		Messages:        messages,
		Store:           false,
		MaxOutputTokens: 10, // minimum is 16
		MaxToolCalls:    1,
		Reasoning: Reasoning{
			Effort: "low",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error mashaling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error when generating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiToken)

	return req, nil
}

// append an input message to a history of messages
func appendHistory(message string, history []Message) []Message {
	messages := []Message{
		{Role: "system", Content: developerPrompt},
	}

	if len(history) > 0 {
		messages = append(messages, history...)
	}

	messages = append(messages, Message{Role: "user", Content: message})

	return messages
}

// parse an OpenAI completion API response
func parseResponse(resp *http.Response) (string, error) {
	var responseData OpenaiResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("error when parsing HTTP response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if responseData.Error != nil {
			return "", fmt.Errorf("openAI API request failed\n  status: [%d]\n  error: [%s]", resp.StatusCode, responseData.Error.Message)
		}
		return "", fmt.Errorf("openAI API request failed\n  status: [%d]", resp.StatusCode)
	}

	if responseData.Error != nil {
		return "", fmt.Errorf("error received with request status %d: [%s]", http.StatusOK, responseData.Error)
	}

	if responseData.Status != "completed" {
		return "", fmt.Errorf("openai API request is not completed: %s", responseData.Status)
	}

	if len(responseData.Output) == 0 {
		return "", fmt.Errorf("no output in response")
	}

	var outputText strings.Builder
	for _, output := range responseData.Output {
		for _, content := range output.Content {
			outputText.WriteString(content.Text)
		}
	}

	return outputText.String(), nil
}
