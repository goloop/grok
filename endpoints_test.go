package grok

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestChatCompletionNative(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &req)
		if string(req.ResponseFormat) != `{"type":"json_object"}` {
			t.Errorf("response_format = %s", req.ResponseFormat)
		}
		io.WriteString(w, `{"model":"m","choices":[{"index":0,`+
			`"message":{"role":"assistant","content":"{}"},"finish_reason":"stop"}]}`)
	})
	defer done()

	resp, err := c.ChatCompletion(context.Background(), &ChatRequest{
		Model:          "m",
		Messages:       []ChatMessage{{Role: "user", Content: "hi"}},
		ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "{}" {
		t.Errorf("content = %v", resp.Choices[0].Message.Content)
	}
}

func TestChatCompletionStreamNative(t *testing.T) {
	events := []string{
		`data: {"choices":[{"index":0,"delta":{"content":"a"}}]}`, ``,
		`data: {"choices":[{"index":0,"delta":{"content":"b"}}]}`, ``,
		`data: [DONE]`, ``,
	}
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		for _, line := range events {
			io.WriteString(w, line+"\n")
		}
	})
	defer done()

	var text strings.Builder
	for chunk, err := range c.ChatCompletionStream(context.Background(), &ChatRequest{
		Model: "m", Messages: []ChatMessage{{Role: "user", Content: "hi"}},
	}) {
		if err != nil {
			t.Fatal(err)
		}
		for _, ch := range chunk.Choices {
			text.WriteString(ch.Delta.Content)
		}
	}
	if text.String() != "ab" {
		t.Errorf("text = %q", text.String())
	}
}

func TestModels(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/models/grok-4"):
			io.WriteString(w, `{"id":"grok-4","object":"model"}`)
		default:
			io.WriteString(w, `{"data":[{"id":"grok-4","object":"model"}]}`)
		}
	})
	defer done()

	ctx := context.Background()
	if models, err := c.Models(ctx); err != nil || len(models) != 1 {
		t.Fatalf("models: %v %+v", err, models)
	}
	m, err := c.GetModel(ctx, "grok-4")
	if err != nil || m.ID != "grok-4" {
		t.Fatalf("get model: %v %+v", err, m)
	}
}

func TestGenerateImage(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/images/generations") {
			t.Errorf("path = %q", r.URL.Path)
		}
		io.WriteString(w, `{"created":1,"data":[{"url":"https://img/x.png"}]}`)
	})
	defer done()

	resp, err := c.GenerateImage(context.Background(), &ImageRequest{
		Model: ModelGrok2Image, Prompt: "a cat", N: 1,
	})
	if err != nil || len(resp.Data) != 1 || resp.Data[0].URL != "https://img/x.png" {
		t.Fatalf("image: %v %+v", err, resp)
	}
}
