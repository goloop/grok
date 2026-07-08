# grok - reference

The full reference for the `grok` package: the client, the shared `goloop/ai`
model, chat completions (interface and native), streaming, image generation and
models.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Creating a client](#creating-a-client)
- [Generate and Stream](#generate-and-stream)
- [Native chat completions](#native-chat-completions)
- [Tools, images and system prompts](#tools-images-and-system-prompts)
- [Image generation](#image-generation)
- [Models](#models)
- [Options and errors](#options-and-errors)

## Mental model

`grok.Client` implements `ai.Client`, the provider-agnostic contract from
`github.com/goloop/ai`. The shared `Generate` and `Stream` cover the common
ground - chat with tools, images and streaming - so code written against the
interface runs on any provider.

Provider-specific power lives in native methods: the full `ChatCompletion`
request, image generation and model listing. Those are not part of the shared
interface. The wire format is chat-completions compatible.

```go
import (
	"github.com/goloop/ai"
	"github.com/goloop/grok"
)
```

## Creating a client

```go
c := grok.New(os.Getenv("XAI_API_KEY"))

c = grok.New(apiKey, grok.WithTimeout(30*time.Second))
```

The base URL defaults to `https://api.x.ai/v1`. Point `WithBaseURL` at any
compatible endpoint to reuse this client against another gateway.

## Generate and Stream

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    grok.ModelGrok4,
	System:   "You are concise.",
	Messages: []ai.Message{ai.UserText("Name three primary colors.")},
})
resp.Text()
resp.ToolCalls()
resp.Usage
```

`Stream` returns `iter.Seq2[ai.Chunk, error]`: text deltas as chunks with
`Text`, a finished tool call as a chunk with `ToolCall`, and a final chunk with
`Done` and `Usage`.

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		return err
	}
	fmt.Print(chunk.Text)
}
```

## Native chat completions

For provider-only options build a `ChatRequest` and call `ChatCompletion` or
`ChatCompletionStream`:

```go
resp, err := c.ChatCompletion(ctx, &grok.ChatRequest{
	Model:          grok.ModelGrok4,
	Messages:       []grok.ChatMessage{{Role: "user", Content: "as JSON"}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

`ChatMessage.Content` is a string or a slice of content parts; `Tools`,
`ToolChoice`, `Temperature`, `TopP`, `MaxTokens`, `Stop`, `N`, `Seed`,
`ResponseFormat` and `User` are all available.

## Tools, images and system prompts

Tool use, images and system prompts use the shared `ai` types: `ai.Tool`,
`ai.Image`, `ai.ToolResult` and a `RoleSystem` message or the `System` field.
Tool results are sent back as `RoleTool` messages whose `ai.ToolResult.ID`
matches the `ai.ToolUse.ID`. Inline image bytes are sent as a base64 data URI.

## Image generation

```go
resp, err := c.GenerateImage(ctx, &grok.ImageRequest{
	Model: grok.ModelGrok2Image, Prompt: "a watercolor cat", N: 1,
})
resp.Data[0].URL // or B64JSON
```

## Models

```go
models, err := c.Models(ctx)
m, err := c.GetModel(ctx, grok.ModelGrok4)
```

## Options and errors

Options: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`,
`WithHeader`.

A non-success response becomes an `*ai.APIError` with `Status`, `Type`, `Code`,
`Message` and the raw body:

```go
var apiErr *ai.APIError
if errors.As(err, &apiErr) && apiErr.Status == http.StatusTooManyRequests {
	// back off
}
```

Requests missing a model or messages fail before the network with
`ai.ErrNoModel` or `ai.ErrNoMessages`.
