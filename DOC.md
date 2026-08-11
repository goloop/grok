# grok - reference

The full reference for the `grok` package: the client, the shared `goloop/ai`
model, chat completions (interface and native), streaming, image generation and
models.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Creating a client](#creating-a-client)
- [Generate and Stream](#generate-and-stream)
- [Structured output](#structured-output)
- [Hosted web search](#hosted-web-search)
- [Capabilities and model-level refusals](#capabilities-and-model-level-refusals)
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

## Structured output

`ai.Request.Format` maps onto the provider's own `response_format`, so a request for JSON
is enforced by the provider rather than merely asked for:

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    "the-model",
	Messages: []ai.Message{ai.UserText("Draft SEO fields for this article.")},
	Format: &ai.Format{
		Type:   ai.FormatJSONSchema,
		Name:   "seo",
		Schema: schema,
	},
})

var seo SEO
err = resp.JSON(&seo)
```

`ai.FormatJSON` goes out as `{"type":"json_object"}` and `ai.FormatJSONSchema`
as `{"type":"json_schema", ...}`. Plain JSON mode also appends
`ai.Format.Instruction()` to the system prompt: this wire format rejects
`json_object` unless the word "json" appears in the messages. Your own system
prompt is kept and the instruction follows it; schema mode leaves it untouched.


`ai.Response.Format` is `ai.FormatNative`: this provider enforces every shape
it accepts. Which models support schema mode is the provider's business - there
is no capability table here, so an unsupported pairing is reported by the
provider itself.

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

## Hosted web search

`ai.Request.Hosted` is answered with `ai.ErrNoHosted` before the request
leaves.

This provider does run searches of its own, but they are not reachable from the
chat endpoint this package speaks, and its documentation leaves more than one
plausible request shape. Declaring support on a guess would mean sending JSON
nobody has seen accepted.

The refusal is the documented behavior, not a gap left in silence: an answer
produced without the search that was asked for looks exactly like one produced
with it, so failing loudly is the only way you can tell them apart. If you would
rather have the answer anyway, ask again without `Hosted`.

## Capabilities and model-level refusals

This driver implements `ai.Capable` and reports no hosted capability - the
same answer `ai.Request.Hosted` already gets here, now available before the
call instead of only as an error after it. A conformance test pins the two
together.

A refusal the provider reports only as a 400 - "this model cannot produce that
format" - now arrives wrapped in `ai.ErrNoFormat`, so one `errors.Is` replaces
matching English prose in an error message; the provider's own `ai.APIError`
stays reachable with `errors.As`. The wrapping is deliberately narrow: only a
400, only a format the request actually asked for, only an error naming that
exact field.

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
