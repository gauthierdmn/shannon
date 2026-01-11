package openai

type Message interface {
	isMessage()
}

type ContentMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (ContentMessage) isMessage() {}

type ToolCallMessage struct {
	Type   string `json:"type"`
	CallId string `json:"call_id"`
	Output string `json:"output"`
}

func (ToolCallMessage) isMessage() {}

type Request struct {
	Model           string    `json:"model"`
	Messages        []Message `json:"input"`
	Store           bool      `json:"store"`
	MaxOutputTokens int       `json:"max_output_tokens"`
	Tools           []*Tool   `json:"tools"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResponseOutput struct {
	Type      string                   `json:"type"`
	Status    string                   `json:"status"`
	Content   []*ResponseOutputContent `json:"content,omitempty"`
	Name      *string                  `json:"name,omitempty"`
	Arguments *string                  `json:"arguments,omitempty"`
	CallId    *string                  `json:"call_id,omitempty"`
}

func (ResponseOutput) isMessage() {}

type ResponseOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Response struct {
	Status string           `json:"status"`
	Error  *ResponseError   `json:"error"`
	Output []ResponseOutput `json:"output"`
}

type Tool struct {
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
