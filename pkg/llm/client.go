package llm

import (
	"context"

	"github.com/gauthierdmn/shannon/pkg/search"
)

type Client interface {
	NewConversation(searchClient search.Client) Conversation
}

type Conversation interface {
	Complete(ctx context.Context, message string) (string, error)
}
