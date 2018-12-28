package storage

import (
	"bytes"
	"context"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine/file"
)

// Storager ...
type Storager interface {
	Close()
	Upload(ctx context.Context, name, contentType string, r io.Reader) error
}

// NewClient ...
func NewClient(ctx context.Context) (Storager, error) {
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

// Client ...
type Client struct {
	*storage.Client
}

// Close ...
func (c *Client) Close() {
	c.Close()
}

// Upload ...
func (c *Client) Upload(ctx context.Context, name, contentType string, r io.Reader) error {
	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		return err
	}

	w := c.Client.Bucket(bucketName).Object(name).NewWriter(ctx)
	defer w.Close()
	w.ContentType = contentType

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	if _, err := io.Copy(w, buf); err != nil {
		return err
	}
	return nil
}
