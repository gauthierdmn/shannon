package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gauthierdmn/shannon/pkg/llm"
	"github.com/gauthierdmn/shannon/pkg/search"
)

const apiEndpoint string = "https://api.openai.com/v1/responses"
const developerPrompt string = "You are a pirate, skip reasoning steps, and answer in less than 10 words."

type Conversation struct {
	client       *Client
	history      []Message
	tools        []*Tool
	searchClient search.Client
}

func (client *Client) NewConversation(searchClient search.Client) llm.Conversation {

	tools := []*Tool{}

	if searchClient != nil {
		tools = append(tools, NewWebSearchTool())
	}

	return &Conversation{
		client:       client,
		history:      nil,
		tools:        tools,
		searchClient: searchClient,
	}
}

// complete a text message using the Open AI completion API
func (conv *Conversation) Complete(ctx context.Context, message string) (string, error) {

	if conv.history == nil {
		conv.history = []Message{
			ContentMessage{Role: "system", Content: developerPrompt},
		}
	}

	conv.history = append(conv.history, ContentMessage{Role: "user", Content: message})

turnLoop:
	for {
		req, err := conv.client.buildRequest(ctx, conv.history, conv.tools)
		if err != nil {
			return "", fmt.Errorf("error when building openai request: %s", err)
		}

		resp, err := conv.client.httpClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("error when sending HTTP request: %w", err)
		}

		defer resp.Body.Close()

		messageResp, err := parseResponse(resp)
		if err != nil {
			return "", fmt.Errorf("error when parsing HTTP response: %w", err)
		}

		for _, output := range messageResp.Output {
			if output.Type == "message" {
				handleMessage(output, conv)
				break turnLoop

			} else if output.Type == "function_call" {
				err := handleFunctionCall(ctx, output, conv)
				if err != nil {
					return "", fmt.Errorf("was not able to run web search")
				}
			}
		}
	}

	lastMessage := conv.history[len(conv.history)-1]

	if contentMsg, ok := lastMessage.(ContentMessage); ok {
		return contentMsg.Content, nil
	}

	return "", fmt.Errorf("bad message type")
}

// build an HTTP request for the OpenAI completion API
func (client *Client) buildRequest(ctx context.Context, messages []Message, tools []*Tool) (*http.Request, error) {

	payload := Request{
		Model:           client.model,
		Messages:        messages,
		Store:           false,
		MaxOutputTokens: 1000, // minimum is 16
		Tools:           tools,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	// useful to debug wrongly formatted payload
	fmt.Println("Req data:", string(jsonData))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error when generating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	return req, nil
}

// parse an OpenAI completion API response
func parseResponse(resp *http.Response) (*Response, error) {
	var responseData Response
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, fmt.Errorf("error when parsing HTTP response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if responseData.Error != nil {
			return nil, fmt.Errorf("openAI API request failed\n  status: [%d]\n  error: [%s]", resp.StatusCode, responseData.Error.Message)
		}
		return nil, fmt.Errorf("openAI API request failed\n  status: [%d]", resp.StatusCode)
	}

	if responseData.Error != nil {
		return nil, fmt.Errorf("error received with request status %d: [%s]", http.StatusOK, responseData.Error)
	}

	if responseData.Status != "completed" {
		return nil, fmt.Errorf("openai API request is not completed: %s", responseData.Status)
	}

	if len(responseData.Output) == 0 {
		return nil, fmt.Errorf("no output in response")
	}

	return &responseData, nil
}

func handleMessage(output ResponseOutput, conv *Conversation) {
	var outputText strings.Builder

	for _, content := range output.Content {
		outputText.WriteString(content.Text)
	}

	conv.history = append(conv.history, ContentMessage{Role: "assistant", Content: outputText.String()})
}

func handleFunctionCall(ctx context.Context, output ResponseOutput, conv *Conversation) error {

	conv.history = append(conv.history, output)

	if output.Name != nil && *output.Name == "web_search" {
		var response SearchToolResponse
		err := json.Unmarshal([]byte(*output.Arguments), &response)
		if err != nil {
			return fmt.Errorf("invalid arguments: %w", err)
		}
		searchResult, err := conv.searchClient.SearchWeb(ctx, response.Query)
		if err != nil {
			return fmt.Errorf("error when running web search")
		}

		conv.history = append(conv.history, ToolCallMessage{Type: "function_call_output", CallId: *output.CallId, Output: searchResult})

		return nil
	}

	return fmt.Errorf("tool with name %s is not supported", *output.Name)
}
