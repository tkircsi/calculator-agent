package agent

import (
	openai "github.com/sashabaranov/go-openai"
)

// ToolHandler wraps an OpenAI tool with its handler function
type ToolHandler struct {
	Tool    openai.Tool
	Handler func(params string) (any, error)
}

// Agent represents our AI agent
type Agent struct {
	client   *openai.Client
	tools    []ToolHandler
	messages []openai.ChatCompletionMessage
}
