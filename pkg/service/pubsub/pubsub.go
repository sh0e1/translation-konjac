package pubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

// Client ...
type Client struct {
	*pubsub.Client
}

// NewClient ...
func NewClient(ctx context.Context, id string) (*Client, error) {
	c, err := pubsub.NewClient(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, err
}

// Close ...
func (c *Client) Close() {
	c.Client.Close()
}

// CreateTopicIfNotExists ...
func (c *Client) CreateTopicIfNotExists(ctx context.Context, id string) (*pubsub.Topic, error) {
	topic := c.Client.Topic(id)
	if ok, err := topic.Exists(ctx); err != nil {
		return nil, err
	} else if !ok {
		if _, err := c.Client.CreateTopic(ctx, id); err != nil {
			return nil, err
		}
	}
	return topic, nil
}

// CreateSubscriptionIfNotExists ...
func (c *Client) CreateSubscriptionIfNotExists(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	sub := c.Client.Subscription(id)
	if ok, err := sub.Exists(ctx); err != nil {
		return nil, err
	} else if !ok {
		if _, err := c.Client.CreateSubscription(ctx, id, cfg); err != nil {
			return nil, err
		}
	}
	return sub, nil
}

// Data ...
type Data struct {
	MessageID      string
	ReplayToken    string
	AudioFilePath  string
	SourceLanguage string
}

// Marshal ...
func (d *Data) Marshal() []byte {
	b, _ := json.Marshal(d)
	return b
}
