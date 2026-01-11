package openai

func NewWebSearchTool() *Tool {
	return &Tool{
		Type:        "function",
		Name:        "web_search",
		Description: "Search the web for content, which is key to answer requests where recent data is critical (weather forecast, news, ...).",
		Parameters: Parameters{
			Type: "object",
			Properties: map[string]Property{
				"query": {
					Type:        "string",
					Description: "The search query to execute",
				},
			},
			Required: []string{"query"},
		},
	}
}

type SearchToolResponse struct {
	Query string `json:"query"`
}
