package webhooks

import (
	"be/internal/apidoc"
	webhookdto "be/internal/modules/webhooks/dto/webhook"
)

// APIDoc documents the webhooks routes for the OpenAPI spec.
func (m *Module) APIDoc() []apidoc.Op {
	const tag = "Webhooks"
	view, modify := resource+":view", resource+":modify"
	return []apidoc.Op{
		{Method: "ANY", Path: "/api/webhooks/capture/:uuid", Tag: tag, Summary: "Capture any HTTP request into the owner's inbox (public)",
			Response: webhookdto.CaptureAckResponse{}, Codes: []string{"not_found", "payload_too_large"}},
		{Method: "GET", Path: "/api/webhooks/inbox", Tag: tag, Summary: "Own inbox (created on first use)", Auth: true, Permission: view, Response: webhookdto.InboxResponse{}},
		{Method: "POST", Path: "/api/webhooks/inbox/regenerate", Tag: tag, Summary: "Rotate the public capture URL", Auth: true, Permission: modify, Response: webhookdto.InboxResponse{}},
		{Method: "GET", Path: "/api/webhooks/inbox/requests", Tag: tag, Summary: "Captured requests", Auth: true, Permission: view, Query: webhookdto.ListRequestsQuery{}, Response: webhookdto.RequestListResponse{}},
		{Method: "GET", Path: "/api/webhooks/inbox/requests/:id", Tag: tag, Summary: "One captured request", Auth: true, Permission: view, Response: webhookdto.RequestDetailResponse{}},
		{Method: "PATCH", Path: "/api/webhooks/inbox/requests/:id/read", Tag: tag, Summary: "Mark read/unread", Auth: true, Permission: view, Body: webhookdto.SetReadRequest{}, Response: webhookdto.RequestDetailResponse{}},
		{Method: "DELETE", Path: "/api/webhooks/inbox/requests", Tag: tag, Summary: "Clear the inbox", Auth: true, Permission: modify, Response: webhookdto.SoftDeleteResponse{}},
		{Method: "DELETE", Path: "/api/webhooks/inbox/requests/:id", Tag: tag, Summary: "Remove one captured request", Auth: true, Permission: modify, Response: webhookdto.SoftDeleteResponse{}},
	}
}
