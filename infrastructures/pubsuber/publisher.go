package pubsuber

import (
	"context"
	"log"
	"saas-billing/app/services"
	"saas-billing/config"

	"cloud.google.com/go/pubsub"
)

type Publisher struct {
	ctx    *context.Context
	client *pubsub.Client
}

func NewPublisher() services.Publisher {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return &Publisher{
		ctx:    &ctx,
		client: client,
	}
}

func (p *Publisher) Publish(topicID string, data []byte) error {
	t := p.client.Topic(topicID)
	result := t.Publish(*p.ctx, &pubsub.Message{
		Data: data,
	})
	_, err := result.Get(*p.ctx)
	if err != nil {
		return err
	}
	return nil
}
