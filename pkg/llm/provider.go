package llm

import (
	"fmt"
	"strings"
)

type Provider string

const (
	ProviderOpenAI Provider = "openai"
)

func ParseProvider(p string) (Provider, error) {
	switch strings.ToLower(p) {
	case string(ProviderOpenAI):
		return ProviderOpenAI, nil
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s (supported: openai)", p)
	}
}
