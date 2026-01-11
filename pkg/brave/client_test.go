package brave

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
	}{
		{
			name:   "empty api key",
			apiKey: "",
		},
		{
			name:   "short api key",
			apiKey: "short",
		},
		{
			name:   "long api key",
			apiKey: "this-is-a-long-api-key-1234567890-;l'l]];[pl]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey)

			if client.apiKey != tt.apiKey {
				t.Fatalf("got key %s, expected %s instead", client.apiKey, tt.apiKey)
			}
		})
	}
}

func TestSearchWeb(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		expectedResult string
	}{
		{
			name: "single result",
			mockResponse: `{
				"web": {
					"results": [
						{
							"age": "1 day",
							"title": "Test",
							"description": "Desc",
							"url": "https://example.com"
						}
					]
				}
			}`,
			expectedResult: `[{"age":"1 day","title":"Test","description":"Desc","url":"https://example.com"}]`,
		},
		{
			name: "multiple results",
			mockResponse: `{
				"web": {
					"results": [
						{
							"age": "1 day",
							"title": "First",
							"description": "First desc",
							"url": "https://example.com/1"
						},
						{
							"age": "2 days",
							"title": "Second",
							"description": "Second desc",
							"url": "https://example.com/2"
						}
					]
				}
			}`,
			expectedResult: `[{"age":"1 day","title":"First","description":"First desc","url":"https://example.com/1"},{"age":"2 days","title":"Second","description":"Second desc","url":"https://example.com/2"}]`,
		},
		{
			name: "empty results",
			mockResponse: `{
				"web": {
					"results": []
				}
			}`,
			expectedResult: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := NewClient("api-key")
			client.apiBaseUrl = server.URL // replace Brave API URL with test server's
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			searchResp, err := client.SearchWeb(ctx, "this is a query on the web!")

			if err != nil {
				t.Fatalf("SearchWeb() error = %v", err)
			}

			if searchResp != tt.expectedResult {
				t.Fatalf("search result is not expected: %s", searchResp)
			}
		})
	}
}
