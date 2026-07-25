package dependency

import (
	webhooksvc "be/internal/services/webhook"
)

func NewWebhookService(infra *Infra) *webhooksvc.Service {
	return webhooksvc.NewService(
		newWebhookInboxRepository(infra.DB),
		newWebhookRequestRepository(infra.DB),
	)
}
