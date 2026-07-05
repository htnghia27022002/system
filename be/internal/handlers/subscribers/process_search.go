package subscribers

import (
	"context"
	"encoding/json"
	"fmt"

	"be/internal/handlers/publisher"
	"be/internal/queue"
	searchsvc "be/internal/services/search"
)

type ProcessSearchOutbox struct {
	processor *searchsvc.IndexProcessor
}

func NewProcessSearchOutbox(processor *searchsvc.IndexProcessor) *ProcessSearchOutbox {
	return &ProcessSearchOutbox{processor: processor}
}

func (h *ProcessSearchOutbox) Name() string {
	return queue.HandlerSearchOutbox
}

func (h *ProcessSearchOutbox) Handle(ctx context.Context, payload []byte) error {
	if h.processor == nil {
		return nil
	}

	var msg publisher.SearchOutboxPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("decode search outbox payload: %w", err)
	}
	if msg.ID == "" {
		return fmt.Errorf("search outbox payload id is required")
	}

	return h.processor.ProcessByID(ctx, msg.ID)
}
