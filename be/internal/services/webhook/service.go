package webhook

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	apperrors "be/internal/common/errors"
	webhookdto "be/internal/dto/webhook"
	webhookmodel "be/internal/models/webhook"
	"be/internal/repository/interfaces"
)

const (
	MaxBodyBytes       = 1 << 20 // 1 MiB
	MaxStoredRequests  = 200
	DefaultListLimit   = 20
	MaxListLimit       = 100
	PublicPathPrefix   = "/tools/webhooks/"
	snippetMaxLen      = 120
)

var (
	// ErrBodyTooLarge is returned when the captured body exceeds 1 MiB.
	ErrBodyTooLarge = errors.New("request body exceeds 1 MiB limit")
)

// CaptureInput is the normalized inbound HTTP capture payload.
type CaptureInput struct {
	Method      string
	URL         string
	ClientIP    string
	Headers     map[string]any
	Query       map[string]any
	Form        map[string]any
	Body        []byte
	ContentType string
	Oversized   bool
}

// Service implements webhook inbox and capture business rules.
type Service struct {
	inboxes  interfaces.WebhookInboxRepository
	requests interfaces.WebhookRequestRepository
}

func NewService(
	inboxes interfaces.WebhookInboxRepository,
	requests interfaces.WebhookRequestRepository,
) *Service {
	return &Service{inboxes: inboxes, requests: requests}
}

func PublicPath(publicUUID string) string {
	return PublicPathPrefix + publicUUID
}

func (s *Service) GetOrCreateInbox(ctx context.Context, userID string) (*webhookdto.InboxResponse, error) {
	if userID == "" {
		return nil, apperrors.ErrUnauthorized
	}
	inbox, err := s.inboxes.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if inbox == nil {
		inbox = &webhookmodel.Inbox{
			UserID:     userID,
			PublicUUID: uuid.NewString(),
		}
		if err := s.inboxes.Create(ctx, inbox); err != nil {
			// Race: another request created the inbox first.
			existing, getErr := s.inboxes.GetByUserID(ctx, userID)
			if getErr != nil {
				return nil, getErr
			}
			if existing != nil {
				inbox = existing
			} else {
				return nil, err
			}
		}
	}
	return toInboxResponse(inbox), nil
}

func (s *Service) RegenerateUUID(ctx context.Context, userID string) (*webhookdto.InboxResponse, error) {
	inbox, err := s.requireOwnerInbox(ctx, userID)
	if err != nil {
		return nil, err
	}
	inbox.PublicUUID = uuid.NewString()
	if err := s.inboxes.Update(ctx, inbox); err != nil {
		return nil, err
	}
	return toInboxResponse(inbox), nil
}

func (s *Service) ListRequests(ctx context.Context, userID string, query webhookdto.ListRequestsQuery) (*webhookdto.RequestListResponse, error) {
	inbox, err := s.requireOwnerInbox(ctx, userID)
	if err != nil {
		return nil, err
	}
	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}

	items, total, err := s.requests.ListActive(ctx, interfaces.WebhookRequestListFilter{
		InboxID:    inbox.ID,
		Method:     query.Method,
		Q:          query.Q,
		ReadFilter: query.Read,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		return nil, err
	}

	out := make([]webhookdto.RequestListItem, 0, len(items))
	for i := range items {
		out = append(out, toListItem(&items[i]))
	}
	return &webhookdto.RequestListResponse{
		Items:           out,
		ActiveCount:     inbox.ActiveCount,
		LifetimeReceived: inbox.LifetimeReceived,
		Page:            page,
		Limit:           limit,
		Total:           total,
		HasMore:         int64(page*limit) < total,
	}, nil
}

func (s *Service) GetRequest(ctx context.Context, userID, requestID string) (*webhookdto.RequestDetailResponse, error) {
	inbox, err := s.requireOwnerInbox(ctx, userID)
	if err != nil {
		return nil, err
	}
	req, err := s.requests.GetByIDAndInbox(ctx, requestID, inbox.ID)
	if err != nil {
		return nil, err
	}
	if req == nil || req.SoftDeletedAt != nil {
		return nil, apperrors.ErrNotFound
	}
	// Opening detail marks the request as read (email-style).
	if !req.IsRead {
		if ok, setErr := s.requests.SetRead(ctx, requestID, inbox.ID, true); setErr != nil {
			return nil, setErr
		} else if ok {
			req.IsRead = true
		}
	}
	return toDetailResponse(req), nil
}

func (s *Service) SetRequestRead(ctx context.Context, userID, requestID string, isRead bool) (*webhookdto.RequestDetailResponse, error) {
	inbox, err := s.requireOwnerInbox(ctx, userID)
	if err != nil {
		return nil, err
	}
	ok, err := s.requests.SetRead(ctx, requestID, inbox.ID, isRead)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	req, err := s.requests.GetByIDAndInbox(ctx, requestID, inbox.ID)
	if err != nil {
		return nil, err
	}
	if req == nil || req.SoftDeletedAt != nil {
		return nil, apperrors.ErrNotFound
	}
	return toDetailResponse(req), nil
}

func (s *Service) SoftDeleteRequest(ctx context.Context, userID, requestID string) (*webhookdto.SoftDeleteResponse, error) {
	if userID == "" {
		return nil, apperrors.ErrUnauthorized
	}
	// Idempotent: no inbox / already soft-deleted / missing ID → 200 ok with current counters.
	inbox, err := s.inboxes.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if inbox == nil {
		return &webhookdto.SoftDeleteResponse{OK: true}, nil
	}
	deleted, err := s.requests.SoftDelete(ctx, requestID, inbox.ID)
	if err != nil {
		return nil, err
	}
	if deleted {
		if inbox.ActiveCount > 0 {
			inbox.ActiveCount--
		}
		if err := s.inboxes.UpdateCounters(ctx, inbox.ID, inbox.LifetimeReceived, inbox.ActiveCount); err != nil {
			return nil, err
		}
	}
	fresh, err := s.inboxes.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if fresh != nil {
		inbox = fresh
	}
	return &webhookdto.SoftDeleteResponse{
		OK:              true,
		ActiveCount:     inbox.ActiveCount,
		LifetimeReceived: inbox.LifetimeReceived,
	}, nil
}

// SoftDeleteAllActive soft-deletes every active request for the caller's inbox.
// Sets activeCount to 0; lifetimeReceived is never decreased. Idempotent when empty / no inbox.
func (s *Service) SoftDeleteAllActive(ctx context.Context, userID string) (*webhookdto.SoftDeleteResponse, error) {
	if userID == "" {
		return nil, apperrors.ErrUnauthorized
	}
	inbox, err := s.inboxes.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if inbox == nil {
		return &webhookdto.SoftDeleteResponse{OK: true}, nil
	}
	if _, err := s.requests.SoftDeleteAllActive(ctx, inbox.ID); err != nil {
		return nil, err
	}
	if err := s.inboxes.UpdateCounters(ctx, inbox.ID, inbox.LifetimeReceived, 0); err != nil {
		return nil, err
	}
	fresh, err := s.inboxes.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if fresh != nil {
		inbox = fresh
	}
	return &webhookdto.SoftDeleteResponse{
		OK:              true,
		ActiveCount:     inbox.ActiveCount,
		LifetimeReceived: inbox.LifetimeReceived,
	}, nil
}

// Capture records an inbound request for the inbox identified by public UUID.
// When oversized is true, stores an error-marker row (no body) and returns ErrBodyTooLarge.
func (s *Service) Capture(ctx context.Context, publicUUID string, in CaptureInput) error {
	if publicUUID == "" {
		return apperrors.ErrNotFound
	}
	inbox, err := s.inboxes.GetByPublicUUID(ctx, publicUUID)
	if err != nil {
		return err
	}
	if inbox == nil {
		return apperrors.ErrNotFound
	}

	body := in.Body
	status := webhookmodel.CaptureStatusOK
	truncated := false
	if in.Oversized {
		body = nil
		status = webhookmodel.CaptureStatusOversized
		truncated = true
	} else if len(body) > MaxBodyBytes {
		body = nil
		status = webhookmodel.CaptureStatusOversized
		truncated = true
		in.Oversized = true
	}

	req := &webhookmodel.Request{
		InboxID:       inbox.ID,
		Method:        strings.ToUpper(strings.TrimSpace(in.Method)),
		URL:           in.URL,
		ClientIP:      in.ClientIP,
		Headers:       normalizeMap(in.Headers),
		Query:         normalizeMap(in.Query),
		Form:          normalizeMap(in.Form),
		Body:          body,
		ContentType:   in.ContentType,
		BodyTruncated: truncated,
		CaptureStatus: status,
	}
	if req.Method == "" {
		req.Method = http.MethodGet
	}

	if err := s.requests.Create(ctx, req); err != nil {
		return err
	}

	inbox.LifetimeReceived++
	inbox.ActiveCount++
	if err := s.inboxes.UpdateCounters(ctx, inbox.ID, inbox.LifetimeReceived, inbox.ActiveCount); err != nil {
		return err
	}

	if err := s.enforceRetention(ctx, inbox); err != nil {
		return err
	}

	if in.Oversized {
		return ErrBodyTooLarge
	}
	return nil
}

// ReadBodyLimited reads at most MaxBodyBytes+1 to detect oversize without unbounded allocation.
func ReadBodyLimited(r io.Reader) (body []byte, oversized bool, err error) {
	if r == nil {
		return nil, false, nil
	}
	limited := io.LimitReader(r, int64(MaxBodyBytes)+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, false, err
	}
	if len(data) > MaxBodyBytes {
		return nil, true, nil
	}
	return data, false, nil
}

func (s *Service) enforceRetention(ctx context.Context, inbox *webhookmodel.Inbox) error {
	count, err := s.requests.CountByInbox(ctx, inbox.ID)
	if err != nil {
		return err
	}
	if count <= MaxStoredRequests {
		return nil
	}
	toDelete := int(count - MaxStoredRequests)
	purged, err := s.requests.DeleteOldest(ctx, inbox.ID, toDelete)
	if err != nil {
		return err
	}
	for _, row := range purged {
		if inbox.LifetimeReceived > 0 {
			inbox.LifetimeReceived--
		}
		if row.SoftDeletedAt == nil && inbox.ActiveCount > 0 {
			inbox.ActiveCount--
		}
	}
	return s.inboxes.UpdateCounters(ctx, inbox.ID, inbox.LifetimeReceived, inbox.ActiveCount)
}

func (s *Service) requireOwnerInbox(ctx context.Context, userID string) (*webhookmodel.Inbox, error) {
	if userID == "" {
		return nil, apperrors.ErrUnauthorized
	}
	inbox, err := s.inboxes.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if inbox == nil {
		return nil, apperrors.ErrNotFound
	}
	return inbox, nil
}

func toInboxResponse(inbox *webhookmodel.Inbox) *webhookdto.InboxResponse {
	return &webhookdto.InboxResponse{
		ID:              inbox.ID,
		PublicUUID:      inbox.PublicUUID,
		PublicPath:      PublicPath(inbox.PublicUUID),
		ActiveCount:     inbox.ActiveCount,
		LifetimeReceived: inbox.LifetimeReceived,
		CreatedAt:       inbox.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       inbox.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toListItem(req *webhookmodel.Request) webhookdto.RequestListItem {
	return webhookdto.RequestListItem{
		ID:        req.ID,
		Method:    req.Method,
		URL:       req.URL,
		ClientIP:  req.ClientIP,
		CreatedAt: req.CreatedAt.UTC().Format(time.RFC3339),
		Snippet:   bodySnippet(req),
		IsRead:    req.IsRead,
	}
}

func toDetailResponse(req *webhookmodel.Request) *webhookdto.RequestDetailResponse {
	body, encoding, isBinary := encodeBodyForJSON(req.Body)
	return &webhookdto.RequestDetailResponse{
		ID:            req.ID,
		InboxID:       req.InboxID,
		Method:        req.Method,
		URL:           req.URL,
		ClientIP:      req.ClientIP,
		Headers:       normalizeMap(req.Headers),
		Query:         normalizeMap(req.Query),
		Form:          normalizeMap(req.Form),
		Body:          body,
		BodyEncoding:  encoding,
		IsBinary:      isBinary,
		ContentType:   req.ContentType,
		BodyTruncated: req.BodyTruncated,
		CaptureStatus: req.CaptureStatus,
		IsRead:        req.IsRead,
		CreatedAt:     req.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func bodySnippet(req *webhookmodel.Request) string {
	if req.CaptureStatus == webhookmodel.CaptureStatusOversized {
		return "[body too large]"
	}
	if len(req.Body) == 0 {
		return ""
	}
	if !utf8.Valid(req.Body) || looksBinary(req.Body) {
		return "[binary]"
	}
	text := string(req.Body)
	if len(text) > snippetMaxLen {
		return text[:snippetMaxLen] + "…"
	}
	return text
}

func encodeBodyForJSON(body []byte) (text, encoding string, isBinary bool) {
	// Contract: bodyEncoding is utf-8 (text, including empty) or base64 (binary).
	if len(body) == 0 {
		return "", "utf-8", false
	}
	if utf8.Valid(body) && !looksBinary(body) {
		return string(body), "utf-8", false
	}
	return base64.StdEncoding.EncodeToString(body), "base64", true
}

func looksBinary(b []byte) bool {
	if bytes.IndexByte(b, 0) >= 0 {
		return true
	}
	return false
}

func normalizeMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}

// Ensure ErrBodyTooLarge is distinguishable from generic errors.
func IsBodyTooLarge(err error) bool {
	return errors.Is(err, ErrBodyTooLarge)
}

// FormatCaptureError keeps a stable message for HTTP clients.
func FormatCaptureError(err error) string {
	if IsBodyTooLarge(err) {
		return "request body exceeds 1 MiB limit"
	}
	return fmt.Sprint(err)
}
