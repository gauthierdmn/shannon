package brave

type WebSearchResponse struct {
	Web WebResult `json:"web"`
}

type WebResult struct {
	Results []WebPageInfo `json:"results"`
}

type WebPageInfo struct {
	Age         string `json:"age"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
