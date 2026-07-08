package grok

import "context"

// ImageRequest is the native image generation request.
type ImageRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"` // "url" or "b64_json"
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
}

// GenerateImage generates one or more images from a prompt.
func (c *Client) GenerateImage(
	ctx context.Context,
	req *ImageRequest,
) (*ImageResponse, error) {
	var out ImageResponse
	if err := c.postJSON(ctx, "/images/generations", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
