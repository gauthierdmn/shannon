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

	"github.com/gauthierdmn/shannon/pkg/llm"
)

func main() {
	modelProvider := flag.String("model-provider", "OpenAI", "The provider of the LLM API.")
	apiToken := flag.String("api-token", "", "A token to access the LLM API.")
	modelName := flag.String("model-name", "gpt-5-nano", "The name of the LLM model to use.")

	flag.Parse()

	if *apiToken == "" {
		log.Fatal("--api-token is required.")
	}

	var client llm.Client

	switch strings.ToLower(*modelProvider) {
	case "openai":
		client = llm.NewOpenaiClient(*modelName, *apiToken)
	default:
		log.Fatalf("LLM API provider [%s] is not supported.", *modelProvider)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conv := client.NewConversation()
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

		answer, err := conv.Complete(ctx, input)

		if err != nil {
			log.Fatalf("Application exited: : %v", err)
		}

		fmt.Println(answer)
	}
}
