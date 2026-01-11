package search

import (
	"fmt"
	"strings"
)

type Provider string

const (
	ProviderBrave Provider = "brave"
)

func ParseProvider(p string) (Provider, error) {
	switch strings.ToLower(p) {
	case string(ProviderBrave):
		return ProviderBrave, nil
	default:
		return "", fmt.Errorf("unsupported search provider: %s (supported: brave)", p)
	}
}
