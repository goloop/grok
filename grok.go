package grok

import "github.com/goloop/ai"

// DefaultBaseURL is the xAI API base URL, including the version segment.
const DefaultBaseURL = "https://api.x.ai/v1"

// Convenience model identifiers. Any model string is accepted; use Models to
// discover what the account can call.
const (
	ModelGrok4       = "grok-4"
	ModelGrok3       = "grok-3"
	ModelGrok3Mini   = "grok-3-mini"
	ModelGrok2Vision = "grok-2-vision-1212"
	ModelGrok2Image  = "grok-2-image-1212"
)

// Client is an xAI (Grok) API client. It implements [ai.Client] and adds the
// provider's native endpoints. The wire format is chat-completions compatible.
type Client struct {
	opts ai.Options
}

var _ ai.Client = (*Client)(nil)

// New returns a Client for the given API key. Shared options (WithBaseURL,
// WithHTTPClient, WithTimeout, WithMaxRetries, WithHeader) configure it.
func New(apiKey string, opts ...Option) *Client {
	s := settings{}
	for _, o := range opts {
		o(&s)
	}

	o := ai.NewOptions(apiKey, s.aiOpts...)
	if o.BaseURL == "" {
		o.BaseURL = DefaultBaseURL
	}

	return &Client{opts: o}
}
