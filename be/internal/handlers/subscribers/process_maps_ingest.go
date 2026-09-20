package subscribers

import (
	"context"
	"encoding/json"
	"fmt"

	"be/internal/handlers/publisher"
	"be/internal/queue"
	ingestsvc "be/internal/services/maps/ingest"
)

type ProcessMapsIngest struct {
	ingest *ingestsvc.Service
}

func NewProcessMapsIngest(ingest *ingestsvc.Service) *ProcessMapsIngest {
	return &ProcessMapsIngest{ingest: ingest}
}

func (h *ProcessMapsIngest) Name() string {
	return queue.HandlerMapsIngest
}

func (h *ProcessMapsIngest) Handle(ctx context.Context, payload []byte) error {
	if h.ingest == nil {
		return nil
	}

	var msg publisher.MapsIngestPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("decode maps ingest payload: %w", err)
	}
	if msg.ID == "" {
		return fmt.Errorf("maps ingest payload id is required")
	}

	return h.ingest.ProcessRun(ctx, msg.ID)
}
