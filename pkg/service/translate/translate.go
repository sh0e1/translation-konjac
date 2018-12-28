package translate

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// Translator ...
type Translator interface {
	Close()
	Translate(ctx context.Context, input, source, target string) (string, error)
	DetectLanguage(ctx context.Context, input string) (string, error)
}

// NewClient ...
func NewClient(ctx context.Context) (Translator, error) {
	c, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

// Client ...
type Client struct {
	*translate.Client
}

// Close ...
func (c *Client) Close() {
	c.Close()
}

// Translate ...
func (c *Client) Translate(ctx context.Context, input, source, target string) (string, error) {
	opt := &translate.Options{
		Format: translate.Text,
	}
	if source != "" {
		opt.Source = language.Make(source)
	}

	translated, err := c.Client.Translate(ctx, []string{input}, language.Make(target), opt)
	if err != nil {
		return "", err
	}
	if len(translated) == 0 {
		return "", fmt.Errorf("have no response translated text: %s", input)
	}
	return strings.TrimRight(translated[0].Text, "\n"), nil
}

// DetectLanguage ...
func (c *Client) DetectLanguage(ctx context.Context, input string) (string, error) {
	detections, err := c.Client.DetectLanguage(ctx, []string{input})
	if err != nil {
		return "", err
	}
	if len(detections) == 0 {
		return "", fmt.Errorf("have no response detected language: %s", input)
	}
	return detections[0][0].Language.String(), nil
}
