package vision

import (
	"context"
	"io"

	"cloud.google.com/go/vision/apiv1"
	"github.com/pkg/errors"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const maxResults = 1

// ImageAnnotator ...
type ImageAnnotator interface {
	Close()
	DetectTexts(ctx context.Context, r io.Reader) (*pb.EntityAnnotation, error)
}

// NewClient ...
func NewClient(ctx context.Context) (ImageAnnotator, error) {
	c, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{ImageAnnotatorClient: c}, nil
}

// Client ...
type Client struct {
	*vision.ImageAnnotatorClient
}

// Close ...
func (c *Client) Close() {
	c.ImageAnnotatorClient.Close()
}

// DetectTexts ...
func (c *Client) DetectTexts(ctx context.Context, r io.Reader) (*pb.EntityAnnotation, error) {
	img, err := vision.NewImageFromReader(r)
	if err != nil {
		return nil, err
	}
	annotations, err := c.ImageAnnotatorClient.DetectTexts(ctx, img, nil, maxResults)
	if err != nil {
		return nil, err
	}
	if len(annotations) == 0 {
		return nil, errors.New("have no response detected image text")
	}
	return annotations[0], nil
}
