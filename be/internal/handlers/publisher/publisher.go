package publisher

import (
	"context"
	"encoding/json"

	"be/internal/queue"
)

type SearchOutboxPayload struct {
	ID string `json:"id"`
}

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

func (p *Publisher) PublishSearchOutbox(ctx context.Context, outboxID string) error {
	if p == nil || p.client == nil || !p.client.Enabled() {
		return nil
	}

	body, err := json.Marshal(SearchOutboxPayload{ID: outboxID})
	if err != nil {
		return err
	}

	return p.client.Publish(ctx, queue.SubjectSearchOutbox, body)
}

func (p *Publisher) Close() {
	if p != nil && p.client != nil {
		p.client.Close()
	}
}
