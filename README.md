[![deps.dev](https://img.shields.io/badge/deps.dev-insights-4c8dbc)](https://deps.dev/go/github.com%2Fgoloop%2Fgrok) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/grok/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/grok) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# grok

`grok` is a Go client for the xAI (Grok) API. It implements the
`github.com/goloop/ai` interface, so it looks and works like every other goloop
AI provider, and exposes the native chat-completions endpoint with its full
options on top.

## Features

- Chat completions: `Generate` for a single response, `Stream` for
  token-by-token output through `iter.Seq2`.
- Tool use (function calling), multimodal image input and system prompts.
- Native `ChatCompletion` and `ChatCompletionStream` with the full option set
  (response_format, seed, n, ...).
- Image generation and model listing.
- Retries on 429 and 5xx with backoff; normalized, typed API errors.
- Depends only on `github.com/goloop/ai` and the standard library.

## Installation

```sh
go get github.com/goloop/grok
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/goloop/ai"
	"github.com/goloop/grok"
)

func main() {
	c := grok.New(os.Getenv("XAI_API_KEY"))

	resp, err := c.Generate(context.Background(), &ai.Request{
		Model:    grok.ModelGrok4,
		Messages: []ai.Message{ai.UserText("Say hello in one word.")},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text())
}
```

## Streaming

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		break
	}
	fmt.Print(chunk.Text)
	if chunk.Done && chunk.Usage != nil {
		fmt.Printf("\n[%d in / %d out]\n",
			chunk.Usage.InputTokens, chunk.Usage.OutputTokens)
	}
}
```

## Tools, images and system prompts

Tools, images and system prompts use the same shared `ai` types as every other
provider (see the [reference](DOC.md)). For provider-only options such as
structured output, build a native `ChatRequest`:

```go
resp, _ := c.ChatCompletion(ctx, &grok.ChatRequest{
	Model:          grok.ModelGrok4,
	Messages:       []grok.ChatMessage{{Role: "user", Content: "List two colors as JSON."}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

## Native endpoints

```go
c.Models(ctx)
c.GetModel(ctx, grok.ModelGrok4)
c.GenerateImage(ctx, &grok.ImageRequest{Model: grok.ModelGrok2Image, Prompt: "a cat"})
```

## Documentation

Full reference: **[DOC.md](DOC.md)** (Ukrainian: **[DOC.UK.md](DOC.UK.md)**).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT - see [LICENSE](LICENSE).
