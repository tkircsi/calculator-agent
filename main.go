package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	ctx := context.Background()

	// Create a reader for stdin
	reader := bufio.NewReader(os.Stdin)

	// Read user input from stdin
	fmt.Print("Enter your message: ")
	userMessage, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	// Trim whitespace and newlines
	userMessage = strings.TrimSpace(userMessage)

	// Process the message and get response
	response, err := agent.ProcessMessage(ctx, userMessage)
	if err != nil {
		log.Fatalf("Error processing message: %v", err)
	}

	// Print the response
	fmt.Printf("Assistant: %s\n", response)
}
