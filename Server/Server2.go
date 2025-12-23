package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetingInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type GreetingOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

type ArithmeticInput struct {
	A float64 `json:"a" jsonschema:"first number"`
	B float64 `json:"b" jsonschema:"second number"`
}

type ArithmeticOutput struct {
	Result float64 `json:"result" jsonschema:"result of the arithmetic operation"`
}

// todo: CallToolRequest에 content 포함하기?

func SayHi(ctx context.Context, req *mcp.CallToolRequest, in GreetingInput) (*mcp.CallToolResult, GreetingOutput, error) {
	return nil, GreetingOutput{Greeting: "Hi " + in.Name}, nil
}

func Add(ctx context.Context, req *mcp.CallToolRequest, in ArithmeticInput) (*mcp.CallToolResult, ArithmeticOutput, error) {
	result := in.A + in.B
	return nil, ArithmeticOutput{Result: result}, nil
}

func Subtract(ctx context.Context, req *mcp.CallToolRequest, in ArithmeticInput) (*mcp.CallToolResult, ArithmeticOutput, error) {
	result := in.A - in.B
	return nil, ArithmeticOutput{Result: result}, nil
}

func Multiply(ctx context.Context, req *mcp.CallToolRequest, in ArithmeticInput) (*mcp.CallToolResult, ArithmeticOutput, error) {
	result := in.A * in.B
	return nil, ArithmeticOutput{Result: result}, nil
}

func Divide(ctx context.Context, req *mcp.CallToolRequest, in ArithmeticInput) (*mcp.CallToolResult, ArithmeticOutput, error) {
	if in.B == 0 {
		return nil, ArithmeticOutput{}, fmt.Errorf("division by zero")
	}
	result := in.A / in.B
	return nil, ArithmeticOutput{Result: result}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi to someone"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "add", Description: "add two numbers"}, Add)
	mcp.AddTool(server, &mcp.Tool{Name: "subtract", Description: "subtract second number from first"}, Subtract)
	mcp.AddTool(server, &mcp.Tool{Name: "multiply", Description: "multiply two numbers"}, Multiply)
	mcp.AddTool(server, &mcp.Tool{Name: "divide", Description: "divide first number by second"}, Divide)

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{},
	)

	http.HandleFunc("/mcp", handler.ServeHTTP)

	log.Println("Starting MCP calculator server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
