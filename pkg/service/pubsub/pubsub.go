package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// Pubsub ...
type Pubsub interface {
	Topic(id string) *pubsub.Topic
	CreateTopic(ctx context.Context, id string) (*pubsub.Topic, error)
}

// Topicker ...
type Topicker interface {
	Exists(ctx context.Context) (bool, error)
	Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult
}

// NewClient ...
func NewClient(id, topicName string) (*pubsub.Client, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, id)
	if err != nil {
		return nil, err
	}

	topic := client.Topic(topicName)
	if ok, err := topic.Exists(ctx); err != nil {
		return nil, err
	} else if !ok {
		if _, err := client.CreateTopic(ctx, topicName); err != nil {
			return nil, err
		}
	}
	return client, err
}
