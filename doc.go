// Package grok is a client for the xAI (Grok) API, built on the goloop/ai
// interface.
//
// The Client implements ai.Client, so Generate and Stream work the same as
// with any other goloop AI provider. On top of that it exposes the native
// chat completions endpoint with its full options, image generation and model
// listing. The wire format is chat-completions compatible.
//
//	c := grok.New(os.Getenv("XAI_API_KEY"))
//	resp, err := c.Generate(ctx, &ai.Request{
//	    Model:    grok.ModelGrok4,
//	    Messages: []ai.Message{ai.UserText("Say hello in one word.")},
//	})
//
// # Structured output
//
// ai.Request.Format maps onto the provider's response_format, so a request for
// JSON is enforced rather than merely asked for, and ai.Response.JSON decodes
// the reply. Plain JSON mode also puts ai.Format.Instruction into the system
// prompt, because this wire format rejects json_object unless the word "json"
// appears in the messages.
//
// # Hosted capabilities
//
// This provider does run searches of its own, but they are not reachable
// from the chat endpoint this package speaks, and its documentation leaves
// more than one plausible request shape. Rather than guess at one, ai.Hosted
// is refused with ai.ErrNoHosted before the request leaves.
//
// The refusal is the documented behavior, not a gap waiting to be filled
// in silence: an answer produced without the search that was asked for
// looks exactly like one produced with it. A caller who would rather have
// the answer anyway asks again without ai.Request.Hosted.
//
// It depends only on goloop/ai and the standard library.
package grok
