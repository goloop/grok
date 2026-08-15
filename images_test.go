package grok

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/goloop/ai"
)

// Bytes decodes an inline image and refuses one delivered as a URL, so a
// caller reads the result the same way it would from any goloop image driver.
func TestImageDataBytes(t *testing.T) {
	inline := ImageData{B64JSON: "aGk="}
	got, err := inline.Bytes()
	if err != nil || string(got) != "hi" {
		t.Fatalf("inline Bytes() = %q, %v", got, err)
	}

	urlOnly := ImageData{URL: "https://img/x.png"}
	if _, err := urlOnly.Bytes(); !errors.Is(err, ErrNoImageBytes) {
		t.Errorf("url Bytes() err = %v, want ErrNoImageBytes", err)
	}

	if _, err := (ImageData{}).Bytes(); !errors.Is(err, ErrNoImageBytes) {
		t.Errorf("empty Bytes() err = %v, want ErrNoImageBytes", err)
	}
}

// A nil request is a clear sentinel, not a panic or an opaque HTTP error.
func TestGenerateImageNilRequest(t *testing.T) {
	c := New("k")
	if _, err := c.GenerateImage(context.Background(), nil); !errors.Is(err, ErrNoImageRequest) {
		t.Errorf("nil request err = %v, want ErrNoImageRequest", err)
	}
}

// Usage, when the provider sends it, is decoded rather than dropped.
func TestImageUsageDecoded(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"created":1,"data":[{"b64_json":"aGk="}],`+
			`"usage":{"total_tokens":42,"input_tokens":10,"output_tokens":32}}`)
	})
	defer done()

	resp, err := c.GenerateImage(context.Background(), &ImageRequest{
		Model: ModelGrok2Image, Prompt: "a cat",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Usage == nil || resp.Usage.TotalTokens != 42 {
		t.Errorf("usage = %+v", resp.Usage)
	}

	png, err := resp.Data[0].Bytes()
	if err != nil || string(png) != "hi" {
		t.Errorf("bytes = %q, %v", png, err)
	}
}

// A response without a usage block leaves Usage nil, distinct from a zero count.
func TestImageUsageAbsent(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"created":1,"data":[{"url":"https://img/x.png"}]}`)
	})
	defer done()

	resp, err := c.GenerateImage(context.Background(), &ImageRequest{
		Model: ModelGrok2Image, Prompt: "a cat",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Usage != nil {
		t.Errorf("usage = %+v, want nil", resp.Usage)
	}
}

// The driver reports image support so a UI can offer the control.
func TestGrokReportsImageSupport(t *testing.T) {
	if !ai.SupportsImages(New("k")) {
		t.Error("grok should report image support")
	}
}
