# Go Agent

A simple Go-based AI agent that can perform mathematical calculations using OpenAI's API with function calling capabilities.

## Features

- Mathematical expression evaluation
- Integration with OpenAI's GPT-4 Turbo
- Interactive command-line interface
- Environment-based configuration

## Agent Capabilities

This agent implements the following AI agent capabilities:

- ✅ **Tool Calling**: The agent can call external tools (calculator) to perform specific tasks
- ✅ **Reasoning**: The agent can understand and process mathematical expressions
- ✅ **Memory**: The agent maintains conversation history
- ✅ **Planning**: The agent can plan the steps needed to solve mathematical problems
- ❌ **Multi-step Reasoning**: The agent currently handles single-step calculations
- ❌ **Self-improvement**: The agent doesn't learn from interactions
- ❌ **Multi-tool Integration**: Currently only implements a calculator tool
- ❌ **Autonomous Operation**: Requires user input for each calculation

## Prerequisites

- Go 1.24 or later
- OpenAI API key

## Setup

1. Clone the repository:
```bash
git clone https://github.com/tkircsi/go-agent.git
cd go-agent
```

2. Create a `.env` file in the project root with your OpenAI API key:
```
OPENAI_API_KEY=your_api_key_here
```

3. Install dependencies:
```bash
go mod download
```

## Usage

Run the agent:
```bash
go run main.go
```

The agent will prompt you to enter a mathematical expression. It will evaluate the expression and return the result.

Example:
```
Enter your message: What is 2 + 2 * 3?
Assistant: The result is 8.00
```

## Project Structure

- `agent/` - Core agent implementation
  - `agent.go` - Main agent logic and OpenAI integration
  - `types.go` - Type definitions
- `main.go` - Entry point and CLI interface

## License

MIT 