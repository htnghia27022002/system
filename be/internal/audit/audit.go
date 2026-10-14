// Package audit records who did what, when, from where. Services call
// Recorder.Record with an action ("user.update") and the affected resource; the
// actor and request facts (IP, user agent, request ID) come from the context,
// filled by middleware.AuditContext and the authenticated route wrapper.
package audit

import (
	"context"
	"log/slog"
	"strings"
	"time"

	auditmodel "be/internal/models/audit"
)

// Entry is one auditable event.
type Entry struct {
	// Action is "<resource>.<verb>", e.g. "user.create", "auth.login_failed".
	Action       string
	ResourceType string
	ResourceID   string
	// Metadata is small, non-secret context (changed field names, counts). Never
	// put passwords, tokens or full request bodies here.
	Metadata map[string]any
	// ActorID / ActorEmail override the request actor, for flows without a JWT
	// yet (sign-in, password reset, accepting an invitation).
	ActorID    string
	ActorEmail string
}

// Recorder persists entries. Implementations must not fail the caller: a
// broken audit sink is logged, the business operation continues.
type Recorder interface {
	Record(ctx context.Context, e Entry)
}

type ctxKey int

const (
	requestKey ctxKey = iota
	actorKey
)

type requestMeta struct {
	IP, UserAgent, RequestID string
}

type actor struct {
	ID, Email string
}

// WithRequest stores request facts for later entries.
func WithRequest(ctx context.Context, ip, userAgent, requestID string) context.Context {
	return context.WithValue(ctx, requestKey, requestMeta{IP: ip, UserAgent: userAgent, RequestID: requestID})
}

// WithActor stores the authenticated user for later entries.
func WithActor(ctx context.Context, id, email string) context.Context {
	return context.WithValue(ctx, actorKey, actor{ID: id, Email: email})
}

// ActorFrom returns the authenticated actor in ctx, if any.
func ActorFrom(ctx context.Context) (id, email string) {
	a, _ := ctx.Value(actorKey).(actor)
	return a.ID, a.Email
}

// Store is the persistence port (implemented by repository.AuditLogRepository).
type Store interface {
	Insert(ctx context.Context, log *auditmodel.Log) error
}

type dbRecorder struct {
	store  Store
	logger *slog.Logger
	now    func() time.Time
}

// NewRecorder writes entries through store. It shares the caller's context, so
// an entry recorded inside a transaction rolls back with it.
func NewRecorder(store Store, logger *slog.Logger) Recorder {
	if logger == nil {
		logger = slog.Default()
	}
	return &dbRecorder{store: store, logger: logger, now: time.Now}
}

func (r *dbRecorder) Record(ctx context.Context, e Entry) {
	row := Build(ctx, e, r.now())
	if err := r.store.Insert(ctx, row); err != nil {
		r.logger.ErrorContext(ctx, "audit record failed", "action", e.Action, "err", err)
	}
}

// Build turns an entry plus context facts into a row (exported for tests and
// alternative sinks).
func Build(ctx context.Context, e Entry, at time.Time) *auditmodel.Log {
	id, email := ActorFrom(ctx)
	if e.ActorID != "" || e.ActorEmail != "" {
		id, email = e.ActorID, e.ActorEmail
	}
	meta, _ := ctx.Value(requestKey).(requestMeta)
	row := &auditmodel.Log{
		ActorEmail:   email,
		Action:       e.Action,
		ResourceType: e.ResourceType,
		ResourceID:   e.ResourceID,
		Metadata:     e.Metadata,
		IPAddress:    meta.IP,
		UserAgent:    clip(meta.UserAgent, 512),
		RequestID:    meta.RequestID,
		CreatedAt:    at.UTC(),
	}
	if id != "" {
		row.ActorID = &id
	}
	if row.Metadata == nil {
		row.Metadata = map[string]any{}
	}
	return row
}

func clip(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) > max {
		return s[:max]
	}
	return s
}

type nop struct{}

func (nop) Record(context.Context, Entry) {}

// Nop discards entries (CLIs, unit tests).
func Nop() Recorder { return nop{} }
