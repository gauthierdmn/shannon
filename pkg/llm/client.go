package llm

import "context"

type Client interface {
	NewConversation() Conversation
}

type Conversation interface {
	Complete(ctx context.Context, message string) (string, error)
}
