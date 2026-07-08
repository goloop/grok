package grok_test

import (
	"encoding/json"
	"fmt"

	"github.com/goloop/ai"
	"github.com/goloop/grok"
)

func ExampleNew() {
	c := grok.New("xai-...")
	_ = c // use c.Generate, c.Stream, c.ChatCompletion, ...
	fmt.Println(grok.ModelGrok4)
	// Output: grok-4
}

// ExampleClient_Generate builds a request. Sending it needs a real API key, so
// this example only shows the shape.
func ExampleClient_Generate() {
	req := &ai.Request{
		Model: grok.ModelGrok4,
		Messages: []ai.Message{
			ai.UserText("Name the capital of France."),
		},
	}
	fmt.Println(req.Model, len(req.Messages))
	// Output: grok-4 1
}

// ExampleTool shows a tool definition passed with a request.
func ExampleTool() {
	tool := ai.Tool{
		Name:        "get_weather",
		Description: "Get the current weather for a city.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"}}}`),
	}
	fmt.Println(tool.Name)
	// Output: get_weather
}
