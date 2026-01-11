package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gauthierdmn/shannon/pkg/brave"
	"github.com/gauthierdmn/shannon/pkg/llm"
	"github.com/gauthierdmn/shannon/pkg/openai"
	"github.com/gauthierdmn/shannon/pkg/search"
)

func main() {
	llmProvider := flag.String("llm-provider", "OpenAI", "The provider of the LLM API.")
	llmApiKey := flag.String("llm-api-key", "", "A token to access the LLM API.")
	modelName := flag.String("model-name", "gpt-5-nano", "The name of the LLM model to use.")

	searchProvider := flag.String("search-provider", "Brave", "The provider of the search API.")
	searchApiKey := flag.String("search-api-key", "", "A token to access the search API.")

	flag.Parse()

	if *llmApiKey == "" {
		log.Fatal("--llm-api-key is required.")
	}

	parsedLlmProvider, err := llm.ParseProvider(*llmProvider)
	if err != nil {
		log.Fatalf("LLM API provider [%s] is not supported.", *llmProvider)
	}

	llmClient, err := NewLlmClient(parsedLlmProvider, *modelName, *llmApiKey)
	if err != nil {
		log.Fatalf("LLM client for provider [%s] was not found.", parsedLlmProvider)
	}

	var searchClient search.Client
	if *searchApiKey != "" {
		parsedSearchProvider, err := search.ParseProvider(*searchProvider)
		if err != nil {
			log.Fatalf("Search API provider [%s] is not supported.", *searchProvider)
		}

		searchClient, err = NewSearchClient(parsedSearchProvider, *searchApiKey)
		if err != nil {
			log.Fatalf("Search client for provider [%s] was not found.", parsedSearchProvider)
		}
	}

	conv := llmClient.NewConversation(searchClient)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println(">")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "" {
			continue
		}

		// defines the maximum time a turn is allowed to take
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		answer, err := conv.Complete(ctx, input)
		cancel()

		if err != nil {
			log.Fatalf("Application exited: : %v", err)
		}

		fmt.Println(answer)
	}
}

func NewLlmClient(provider llm.Provider, model, apiKey string) (llm.Client, error) {
	switch provider {
	case llm.ProviderOpenAI:
		return openai.NewClient(model, apiKey), nil
	default:
		return nil, fmt.Errorf("provider %s has no client implementation", provider)
	}
}

func NewSearchClient(provider search.Provider, apiKey string) (search.Client, error) {
	switch provider {
	case search.ProviderBrave:
		return brave.NewClient(apiKey), nil
	default:
		return nil, fmt.Errorf("provider %s has no client implementation", provider)
	}
}
