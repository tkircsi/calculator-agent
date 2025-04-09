package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tkircsi/go-agent/agent"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
		os.Exit(1)
	}

	// Create a new agent
	agent := agent.NewAgent(apiKey)

	// Example conversation
	ctx := context.Background()

	// Test messages
	messages := []string{
		"Hello! Can you help me calculate 2 + 2?",
		// "What about 5 * 5?",
		// "Tell me about yourself",
		"Calculate 10 / 5 and 20 / 4",
	}

	for _, msg := range messages {
		fmt.Printf("\nUser: %s\n", msg)
		response, err := agent.ProcessMessage(ctx, msg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Assistant: %s\n", response)
	}
}
