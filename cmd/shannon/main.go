package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gauthierdmn/shannon/pkg/llm"
)

func main() {
	apiToken := flag.String("api-token", "", "Token to access LLM API.")
	message := flag.String("message", "", "Message to send to the LLM.")

	flag.Parse()

	if *apiToken == "" {
		log.Fatal("--api-token is required.")
	}
	if *message == "" {
		log.Fatal("--message is required")
	}

	client := llm.New("gpt-5-nano", *apiToken)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	answer, err := client.Complete(ctx, *message, nil)

	if err != nil {
		log.Fatalf("Application exited: : %v", err)
	}

	fmt.Println("LLM answer:", answer)
}
