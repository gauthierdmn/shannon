package main

import (
	"flag"
	"fmt"
	"log"

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

	config := llm.New("gpt-5-nano", *apiToken)

	answer, err := config.Complete(*message, nil)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("LLM answer:", answer)
}
