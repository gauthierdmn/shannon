package search

import "context"

type Client interface {
	SearchWeb(ctx context.Context, query string) (string, error)
}
