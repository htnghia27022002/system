// Package auditlog serves the admin view of the audit trail.
package auditlog

import (
	"context"
	"fmt"
	"strings"
	"time"

	apperrors "be/common/errors"
	"be/common/utils"
	admindto "be/internal/dto/admin"
	auditmodel "be/internal/models/audit"
	"be/internal/repository/interfaces"
	"be/pkg/query"
)

type Service struct {
	repo interfaces.AuditLogRepository
}

func NewService(repo interfaces.AuditLogRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, form admindto.ListAuditLogsQuery) (query.Page[admindto.AuditLogResponse], error) {
	page, pageSize, _, _ := form.Pagination()
	q := query.New(page, pageSize).
		OrderBy("created_at DESC").
		WhereEqual("action", strings.TrimSpace(form.Action)).
		WhereEqual("actor_id", form.ActorID).
		WhereEqual("resource_type", strings.TrimSpace(form.ResourceType)).
		WhereEqual("resource_id", strings.TrimSpace(form.ResourceID)).
		WhereLikeAny([]string{"action", "actor_email", "resource_id"}, strings.TrimSpace(form.Search))
	if t, err := parseTime(form.From); err != nil {
		return query.Page[admindto.AuditLogResponse]{}, err
	} else if t != nil {
		q.WhereRaw("created_at >= ?", *t)
	}
	if t, err := parseTime(form.To); err != nil {
		return query.Page[admindto.AuditLogResponse]{}, err
	} else if t != nil {
		q.WhereRaw("created_at < ?", *t)
	}

	rows, total, err := s.repo.List(ctx, q)
	if err != nil {
		return query.Page[admindto.AuditLogResponse]{}, err
	}
	items := make([]admindto.AuditLogResponse, 0, len(rows))
	for i := range rows {
		items = append(items, toResponse(&rows[i]))
	}
	return query.NewPage(items, total, q.Page, q.PageSize), nil
}

func toResponse(row *auditmodel.Log) admindto.AuditLogResponse {
	meta := row.Metadata
	if meta == nil {
		meta = map[string]any{}
	}
	return admindto.AuditLogResponse{
		ID:           row.ID,
		ActorID:      row.ActorID,
		ActorEmail:   row.ActorEmail,
		Action:       row.Action,
		ResourceType: row.ResourceType,
		ResourceID:   row.ResourceID,
		Metadata:     meta,
		IPAddress:    row.IPAddress,
		UserAgent:    row.UserAgent,
		RequestID:    row.RequestID,
		CreatedAt:    utils.FormatRFC3339(row.CreatedAt),
	}
}

func parseTime(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf("%w: dates must be RFC 3339 (e.g. 2026-01-31T00:00:00Z)", apperrors.ErrBadRequest)
	}
	return &t, nil
}
