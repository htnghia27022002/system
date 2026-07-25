package testutil

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	webhookmodel "be/internal/models/webhook"
	"be/internal/repository/interfaces"
)

// MemoryWebhookInboxRepo is an in-memory WebhookInboxRepository for unit tests.
type MemoryWebhookInboxRepo struct {
	mu      sync.Mutex
	ByID    map[string]*webhookmodel.Inbox
	ByUser  map[string]string
	ByUUID  map[string]string
}

func NewMemoryWebhookInboxRepo() *MemoryWebhookInboxRepo {
	return &MemoryWebhookInboxRepo{
		ByID:   map[string]*webhookmodel.Inbox{},
		ByUser: map[string]string{},
		ByUUID: map[string]string{},
	}
}

func (m *MemoryWebhookInboxRepo) GetByUserID(_ context.Context, userID string) (*webhookmodel.Inbox, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.ByUser[userID]
	if !ok {
		return nil, nil
	}
	return cloneInbox(m.ByID[id]), nil
}

func (m *MemoryWebhookInboxRepo) GetByPublicUUID(_ context.Context, publicUUID string) (*webhookmodel.Inbox, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.ByUUID[publicUUID]
	if !ok {
		return nil, nil
	}
	return cloneInbox(m.ByID[id]), nil
}

func (m *MemoryWebhookInboxRepo) Create(_ context.Context, inbox *webhookmodel.Inbox) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if inbox.ID == "" {
		inbox.ID = uuid.NewString()
	}
	if inbox.PublicUUID == "" {
		inbox.PublicUUID = uuid.NewString()
	}
	now := time.Now().UTC()
	if inbox.CreatedAt.IsZero() {
		inbox.CreatedAt = now
	}
	inbox.UpdatedAt = now
	cp := cloneInbox(inbox)
	m.ByID[cp.ID] = cp
	m.ByUser[cp.UserID] = cp.ID
	m.ByUUID[cp.PublicUUID] = cp.ID
	*inbox = *cloneInbox(cp)
	return nil
}

func (m *MemoryWebhookInboxRepo) Update(_ context.Context, inbox *webhookmodel.Inbox) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	existing, ok := m.ByID[inbox.ID]
	if !ok {
		return nil
	}
	delete(m.ByUUID, existing.PublicUUID)
	inbox.UpdatedAt = time.Now().UTC()
	cp := cloneInbox(inbox)
	m.ByID[cp.ID] = cp
	m.ByUser[cp.UserID] = cp.ID
	m.ByUUID[cp.PublicUUID] = cp.ID
	return nil
}

func (m *MemoryWebhookInboxRepo) UpdateCounters(_ context.Context, inboxID string, lifetimeReceived, activeCount int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	inbox, ok := m.ByID[inboxID]
	if !ok {
		return nil
	}
	inbox.LifetimeReceived = lifetimeReceived
	inbox.ActiveCount = activeCount
	inbox.UpdatedAt = time.Now().UTC()
	return nil
}

// MemoryWebhookRequestRepo is an in-memory WebhookRequestRepository for unit tests.
type MemoryWebhookRequestRepo struct {
	mu   sync.Mutex
	ByID map[string]*webhookmodel.Request
}

func NewMemoryWebhookRequestRepo() *MemoryWebhookRequestRepo {
	return &MemoryWebhookRequestRepo{ByID: map[string]*webhookmodel.Request{}}
}

func (m *MemoryWebhookRequestRepo) Create(_ context.Context, req *webhookmodel.Request) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if req.ID == "" {
		req.ID = uuid.NewString()
	}
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now().UTC()
	}
	if req.Headers == nil {
		req.Headers = map[string]any{}
	}
	if req.Query == nil {
		req.Query = map[string]any{}
	}
	if req.Form == nil {
		req.Form = map[string]any{}
	}
	m.ByID[req.ID] = cloneRequest(req)
	return nil
}

func (m *MemoryWebhookRequestRepo) GetByIDAndInbox(_ context.Context, id, inboxID string) (*webhookmodel.Request, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.ByID[id]
	if !ok || req.InboxID != inboxID {
		return nil, nil
	}
	return cloneRequest(req), nil
}

func (m *MemoryWebhookRequestRepo) ListActive(_ context.Context, filter interfaces.WebhookRequestListFilter) ([]webhookmodel.Request, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var items []webhookmodel.Request
	for _, req := range m.ByID {
		if req.InboxID != filter.InboxID || req.SoftDeletedAt != nil {
			continue
		}
		if method := strings.TrimSpace(filter.Method); method != "" && !strings.EqualFold(req.Method, method) {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(filter.ReadFilter)) {
		case "read":
			if !req.IsRead {
				continue
			}
		case "unread":
			if req.IsRead {
				continue
			}
		}
		if term := strings.TrimSpace(filter.Q); term != "" {
			if !strings.Contains(strings.ToLower(req.URL), strings.ToLower(term)) &&
				!strings.Contains(strings.ToLower(string(req.Body)), strings.ToLower(term)) {
				continue
			}
		}
		items = append(items, *cloneRequest(req))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	total := int64(len(items))
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	start := (page - 1) * limit
	if start >= len(items) {
		return []webhookmodel.Request{}, total, nil
	}
	end := start + limit
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], total, nil
}

func (m *MemoryWebhookRequestRepo) SoftDelete(_ context.Context, id, inboxID string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.ByID[id]
	if !ok || req.InboxID != inboxID || req.SoftDeletedAt != nil {
		return false, nil
	}
	now := time.Now().UTC()
	req.SoftDeletedAt = &now
	return true, nil
}

func (m *MemoryWebhookRequestRepo) SoftDeleteAllActive(_ context.Context, inboxID string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	var n int64
	for _, req := range m.ByID {
		if req.InboxID != inboxID || req.SoftDeletedAt != nil {
			continue
		}
		req.SoftDeletedAt = &now
		n++
	}
	return n, nil
}

func (m *MemoryWebhookRequestRepo) SetRead(_ context.Context, id, inboxID string, isRead bool) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.ByID[id]
	if !ok || req.InboxID != inboxID || req.SoftDeletedAt != nil {
		return false, nil
	}
	req.IsRead = isRead
	return true, nil
}

func (m *MemoryWebhookRequestRepo) CountByInbox(_ context.Context, inboxID string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var count int64
	for _, req := range m.ByID {
		if req.InboxID == inboxID {
			count++
		}
	}
	return count, nil
}

func (m *MemoryWebhookRequestRepo) DeleteOldest(_ context.Context, inboxID string, limit int) ([]webhookmodel.Request, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if limit <= 0 {
		return nil, nil
	}
	var items []*webhookmodel.Request
	for _, req := range m.ByID {
		if req.InboxID == inboxID {
			items = append(items, req)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	if len(items) > limit {
		items = items[:limit]
	}
	out := make([]webhookmodel.Request, 0, len(items))
	for _, req := range items {
		out = append(out, *cloneRequest(req))
		delete(m.ByID, req.ID)
	}
	return out, nil
}

func cloneInbox(in *webhookmodel.Inbox) *webhookmodel.Inbox {
	if in == nil {
		return nil
	}
	cp := *in
	return &cp
}

func cloneRequest(in *webhookmodel.Request) *webhookmodel.Request {
	if in == nil {
		return nil
	}
	cp := *in
	if in.Body != nil {
		cp.Body = append([]byte(nil), in.Body...)
	}
	if in.SoftDeletedAt != nil {
		t := *in.SoftDeletedAt
		cp.SoftDeletedAt = &t
	}
	cp.Headers = copyAnyMap(in.Headers)
	cp.Query = copyAnyMap(in.Query)
	cp.Form = copyAnyMap(in.Form)
	return &cp
}

func copyAnyMap(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
