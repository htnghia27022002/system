// Package publisher is the outbound side of the queue: modules publish JSON
// payloads to their own subjects through it. It is a no-op when NATS is disabled.
package publisher

import (
	"context"
	"encoding/json"

	"be/internal/queue"
)

type Publisher struct {
	client *queue.Client
}

func New(client *queue.Client) *Publisher {
	return &Publisher{client: client}
}

func NewNoop() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Enabled() bool {
	return p != nil && p.client != nil && p.client.Enabled()
}

// PublishJSON marshals payload and publishes it to subject. No-op when disabled.
func (p *Publisher) PublishJSON(ctx context.Context, subject string, payload any) error {
	if !p.Enabled() {
		return nil
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.client.Publish(ctx, subject, body)
}

func (p *Publisher) Close() {
	if p != nil && p.client != nil {
		p.client.Close()
	}
}
