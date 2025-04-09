package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Knetic/govaluate"
	openai "github.com/sashabaranov/go-openai"
)

// NewAgent creates a new agent instance
func NewAgent(apiKey string) *Agent {
	client := openai.NewClient(apiKey)

	// Initialize tools
	tools := []ToolHandler{
		{
			Tool: openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        "calculator",
					Description: "Calculate the result of a mathematical expression",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"expression": map[string]any{
								"type":        "string",
								"description": "The mathematical expression to calculate (e.g., '2 + 2', '10 / 2 * 5')",
							},
						},
						"required": []string{"expression"},
					},
				},
			},
			Handler: func(params string) (any, error) {
				// Parse the input as JSON
				var input struct {
					Expression string `json:"expression"`
				}
				if err := json.Unmarshal([]byte(params), &input); err != nil {
					return nil, err
				}

				// Use Go's expression evaluator
				expr, err := govaluate.NewEvaluableExpression(input.Expression)
				if err != nil {
					return nil, fmt.Errorf("invalid expression: %s", input.Expression)
				}

				result, err := expr.Evaluate(nil)
				if err != nil {
					return nil, fmt.Errorf("error evaluating expression: %v", err)
				}

				// Format the result
				var formattedResult string
				switch v := result.(type) {
				case float64:
					formattedResult = fmt.Sprintf("%.2f", v)
				case int:
					formattedResult = fmt.Sprintf("%d", v)
				default:
					formattedResult = fmt.Sprintf("%v", v)
				}

				return map[string]any{
					"result": formattedResult,
				}, nil
			},
		},
	}

	// Initialize messages with system message
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are a calculator agent. Your task is to:
1. If the user asks for a calculation or there is a mathematical expression to evaluate, use the calculator tool to perform the calculation.
2. If the user's message doesn't contain any calculation or mathematical expression, respond with "There is nothing to calculate in your message."
3. Keep your responses focused on calculations and mathematical operations.`,
		},
	}

	return &Agent{
		client:   client,
		tools:    tools,
		messages: messages,
	}
}

// ProcessMessage handles a user message and returns the agent's response
func (a *Agent) ProcessMessage(ctx context.Context, userMessage string) (string, error) {
	// Add user message to conversation history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	// Prepare tools for the API call
	openaiTools := make([]openai.Tool, len(a.tools))
	for i, tool := range a.tools {
		openaiTools[i] = tool.Tool
	}

	// Create the completion request with tools
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    "gpt-4-turbo-preview",
			Messages: a.messages,
			Tools:    openaiTools,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error creating completion: %v", err)
	}

	// Get the assistant's response
	assistantMessage := resp.Choices[0].Message

	// Check if the model wants to call any functions
	if assistantMessage.ToolCalls != nil && len(assistantMessage.ToolCalls) > 0 {
		log.Print("Function calls detected")

		// Add the assistant's message with tool calls to the conversation
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:      openai.ChatMessageRoleAssistant,
			Content:   assistantMessage.Content,
			ToolCalls: assistantMessage.ToolCalls,
		})

		// Process each tool call
		for _, toolCall := range assistantMessage.ToolCalls {
			log.Printf("Processing tool call: %s", toolCall.Function.Name)

			// Find the tool that matches the function call
			for _, tool := range a.tools {
				if tool.Tool.Function.Name == toolCall.Function.Name {
					log.Printf("Tool found: %s", tool.Tool.Function.Name)

					// Call the tool with the function arguments
					result, err := tool.Handler(toolCall.Function.Arguments)
					if err != nil {
						return "", fmt.Errorf("tool error: %v", err)
					}

					// Convert result to JSON
					resultJSON, err := json.Marshal(result)
					if err != nil {
						return "", fmt.Errorf("error marshaling result: %v", err)
					}

					// Add the tool result to the conversation
					a.messages = append(a.messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    string(resultJSON),
						ToolCallID: toolCall.ID,
					})
				}
			}
		}

		// Get the final response from the model after all tool calls
		finalResp, err := a.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    "gpt-4-turbo-preview",
				Messages: a.messages,
			},
		)
		if err != nil {
			return "", fmt.Errorf("error creating final completion: %v", err)
		}

		return finalResp.Choices[0].Message.Content, nil
	}

	// If no function was called, just return the assistant's response
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: assistantMessage.Content,
	})

	return assistantMessage.Content, nil
}
