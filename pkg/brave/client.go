package brave

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const searchApiUrl string = "https://api.search.brave.com/res/v1/web/search"
const searchResultCount int = 3

type Client struct {
	httpClient *http.Client
	apiBaseUrl string
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		apiBaseUrl: searchApiUrl,
		apiKey:     apiKey,
	}
}

func (client *Client) SearchWeb(ctx context.Context, query string) (string, error) {

	params := url.Values{}
	params.Add("q", query)
	params.Add("count", strconv.Itoa(searchResultCount))

	searchUrl := client.apiBaseUrl + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error when generating HTTP request: %w", err)
	}

	req.Header.Set("X-Subscription-Token", client.apiKey)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error when sending HTTP request: %w", err)
	}

	defer resp.Body.Close()

	searchResp, err := parseResponse(resp)
	if err != nil {
		return "", fmt.Errorf("error when parsing HTTP response: %w", err)
	}

	fmt.Println("Search result:", searchResp)

	return searchResp, nil
}

func parseResponse(resp *http.Response) (string, error) {
	var responseData WebSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("error when parsing HTTP response: %w", err)
	}

	jsonBytes, err := json.Marshal(responseData.Web.Results)
	if err != nil {
		return "", fmt.Errorf("error when converting web search result to JSON string")
	}

	return string(jsonBytes), nil
}
