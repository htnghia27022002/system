package testutil

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"be/internal/mail"
	accountmodel "be/internal/models/account"
	"be/internal/repository/interfaces"
)

var _ interfaces.UserTokenRepository = (*MockUserTokenRepo)(nil)

// MockUserTokenRepo is an in-memory UserTokenRepository.
type MockUserTokenRepo struct {
	Tokens map[string]*accountmodel.UserToken // by id
}

func NewMockUserTokenRepo() *MockUserTokenRepo {
	return &MockUserTokenRepo{Tokens: map[string]*accountmodel.UserToken{}}
}

func (m *MockUserTokenRepo) Create(_ context.Context, token *accountmodel.UserToken) error {
	if token.ID == "" {
		token.ID = uuid.NewString()
	}
	copied := *token
	m.Tokens[token.ID] = &copied
	return nil
}

func (m *MockUserTokenRepo) FindValid(_ context.Context, purpose accountmodel.TokenPurpose, tokenHash string) (*accountmodel.UserToken, error) {
	for _, t := range m.Tokens {
		if t.TokenHash == tokenHash && t.Purpose == purpose && t.UsedAt == nil && t.ExpiresAt.After(time.Now()) {
			copied := *t
			return &copied, nil
		}
	}
	return nil, nil
}

func (m *MockUserTokenRepo) Consume(_ context.Context, id string) (bool, error) {
	t, ok := m.Tokens[id]
	if !ok || t.UsedAt != nil {
		return false, nil
	}
	now := time.Now()
	t.UsedAt = &now
	return true, nil
}

func (m *MockUserTokenRepo) InvalidateOpen(_ context.Context, userID string, purpose accountmodel.TokenPurpose) error {
	now := time.Now()
	for _, t := range m.Tokens {
		if t.UserID == userID && t.Purpose == purpose && t.UsedAt == nil {
			t.UsedAt = &now
		}
	}
	return nil
}

func (m *MockUserTokenRepo) LatestCreatedAt(_ context.Context, userID string, purpose accountmodel.TokenPurpose) (*time.Time, error) {
	var latest *time.Time
	for _, t := range m.Tokens {
		if t.UserID == userID && t.Purpose == purpose && t.UsedAt == nil {
			if latest == nil || t.CreatedAt.After(*latest) {
				created := t.CreatedAt
				latest = &created
			}
		}
	}
	return latest, nil
}

// NoTx runs fn directly (unit tests have no database).
type NoTx struct{}

func (NoTx) WithTx(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) }

// RecordingMailer captures sent messages instead of delivering them.
type RecordingMailer struct {
	mu   sync.Mutex
	Sent []mail.Message
}

func (r *RecordingMailer) Send(_ context.Context, msg mail.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Sent = append(r.Sent, msg)
	return nil
}

// Last returns the most recent message, or the zero Message.
func (r *RecordingMailer) Last() mail.Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Sent) == 0 {
		return mail.Message{}
	}
	return r.Sent[len(r.Sent)-1]
}

// Count returns how many messages were sent.
func (r *RecordingMailer) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.Sent)
}
