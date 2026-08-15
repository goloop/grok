package grok

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
)

// The values ImageRequest.ResponseFormat accepts. The model constant
// (ModelGrok2Image) lives in grok.go with the other model identifiers.
const (
	ImageFormatURL     = "url"
	ImageFormatB64JSON = "b64_json"
)

// Errors reported for an image request or its result. They carry the same
// names as in the other goloop drivers, so application code reads the same way
// whichever provider is behind it.
var (
	// ErrNoImageRequest is returned by GenerateImage for a nil request.
	ErrNoImageRequest = errors.New("grok: image request is nil")

	// ErrNoImageBytes is returned by ImageData.Bytes for an image the provider
	// returned as a URL. Fetch the URL, or ask for ImageFormatB64JSON.
	ErrNoImageBytes = errors.New("grok: image was returned as a URL")
)

// ImageRequest is the native image generation request.
//
// The xAI images endpoint takes the OpenAI shape but a smaller subset of it:
// there is no size, quality or style knob, so those fields are deliberately
// absent rather than accepted and ignored. ResponseFormat is kept because xAI
// does honor it.
type ImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n,omitempty"`

	// ResponseFormat is ImageFormatURL or ImageFormatB64JSON.
	ResponseFormat string `json:"response_format,omitempty"`
}

// ImageData is one generated image: a URL or base64-encoded bytes.
type ImageData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// ImageResponse is the native image generation response.
type ImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`

	// Usage carries the tokens the request billed, when the provider reports
	// them. It is a pointer so a nil Usage - the provider said nothing - is
	// distinct from a zero count. The type mirrors the other drivers'.
	Usage *ImageUsage `json:"usage,omitempty"`
}

// ImageUsage reports the tokens an image request consumed. It matches the
// shape used by the other goloop image drivers so accounting code is uniform;
// fields the provider does not send stay zero.
type ImageUsage struct {
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Bytes returns the decoded image. It reads what the provider put inline and
// does no I/O: an image delivered as a URL returns [ErrNoImageBytes], because
// fetching it is a network call with the caller's own timeouts, redirects and
// proxy rules, not something a field accessor should decide.
func (d ImageData) Bytes() ([]byte, error) {
	if d.B64JSON == "" {
		if d.URL != "" {
			return nil, fmt.Errorf("%w: %s", ErrNoImageBytes, d.URL)
		}
		return nil, ErrNoImageBytes
	}
	return base64.StdEncoding.DecodeString(d.B64JSON)
}

// GenerateImage generates one or more images from a prompt. Read each image
// with [ImageData.Bytes].
func (c *Client) GenerateImage(
	ctx context.Context,
	req *ImageRequest,
) (*ImageResponse, error) {
	if req == nil {
		return nil, ErrNoImageRequest
	}
	var out ImageResponse
	if err := c.postJSON(ctx, "/images/generations", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
