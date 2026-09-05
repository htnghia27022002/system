package testutil

import (
	"context"
	"time"

	"github.com/google/uuid"

	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
	"be/pkg/query"
)

var _ interfaces.AuthRepository = (*MockAuthRepo)(nil)

// MockAuthRepo is an in-memory AuthRepository for unit tests.
type MockAuthRepo struct {
	RefreshTokens map[string]*authmodel.RefreshToken
}

func refreshTokenActive(token *authmodel.RefreshToken) bool {
	if token == nil || token.RevokedAt != nil {
		return false
	}
	if !token.ExpiresAt.IsZero() && !token.ExpiresAt.After(time.Now()) {
		return false
	}
	return true
}

func (m *MockAuthRepo) CreateRefreshToken(_ context.Context, token *authmodel.RefreshToken) error {
	if m.RefreshTokens == nil {
		m.RefreshTokens = map[string]*authmodel.RefreshToken{}
	}
	if token.ID == "" {
		token.ID = uuid.NewString()
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}
	if token.LastUsedAt.IsZero() {
		token.LastUsedAt = time.Now()
	}
	m.RefreshTokens[token.TokenHash] = token
	return nil
}

func (m *MockAuthRepo) FindRefreshTokenByHash(_ context.Context, tokenHash string) (*authmodel.RefreshToken, error) {
	token, ok := m.RefreshTokens[tokenHash]
	if !ok || !refreshTokenActive(token) {
		return nil, nil
	}
	return token, nil
}

func (m *MockAuthRepo) RevokeRefreshToken(_ context.Context, tokenHash string, revokedAt time.Time) error {
	if token, ok := m.RefreshTokens[tokenHash]; ok {
		token.RevokedAt = &revokedAt
	}
	return nil
}

func (m *MockAuthRepo) ListActiveByUserID(_ context.Context, userID string) ([]authmodel.RefreshToken, error) {
	out := make([]authmodel.RefreshToken, 0)
	for _, token := range m.RefreshTokens {
		if token.UserID != userID || !refreshTokenActive(token) {
			continue
		}
		out = append(out, *token)
	}
	return out, nil
}

func (m *MockAuthRepo) RevokeByID(_ context.Context, userID, id string) (int64, error) {
	now := time.Now()
	for _, token := range m.RefreshTokens {
		if token.ID != id || token.UserID != userID || !refreshTokenActive(token) {
			continue
		}
		token.RevokedAt = &now
		return 1, nil
	}
	return 0, nil
}

func (m *MockAuthRepo) RevokeAllExcept(_ context.Context, userID, keepID string) (int64, error) {
	now := time.Now()
	var n int64
	for _, token := range m.RefreshTokens {
		if token.UserID != userID || token.ID == keepID || !refreshTokenActive(token) {
			continue
		}
		token.RevokedAt = &now
		n++
	}
	return n, nil
}

func (m *MockAuthRepo) FindOAuthAccount(context.Context, string, string) (*authmodel.OAuthAccount, error) {
	return nil, nil
}

func (m *MockAuthRepo) CreateOAuthAccount(context.Context, *authmodel.OAuthAccount) error {
	return nil
}

func (m *MockAuthRepo) UpdateOAuthAccount(context.Context, *authmodel.OAuthAccount) error {
	return nil
}

func (m *MockAuthRepo) ListProvidersByUserIDs(context.Context, []string) (map[string][]string, error) {
	return map[string][]string{}, nil
}

// MockRoleRepo is an in-memory RoleRepository for unit tests.
type MockRoleRepo struct {
	Roles       map[string]*rolemodel.Role
	Permissions map[string][]string
}

func (m *MockRoleRepo) Create(context.Context, *rolemodel.Role) error { return nil }
func (m *MockRoleRepo) GetByID(_ context.Context, id string) (*rolemodel.Role, error) {
	return m.Roles[id], nil
}
func (m *MockRoleRepo) GetBySlug(_ context.Context, slug string) (*rolemodel.Role, error) {
	for _, role := range m.Roles {
		if role.Slug == slug {
			return role, nil
		}
	}
	return nil, nil
}
func (m *MockRoleRepo) ListAll(context.Context) ([]rolemodel.Role, error) {
	return nil, nil
}
func (m *MockRoleRepo) List(context.Context, *query.Query) ([]rolemodel.Role, int64, error) {
	return nil, 0, nil
}
func (m *MockRoleRepo) Update(context.Context, *rolemodel.Role) error { return nil }
func (m *MockRoleRepo) Delete(context.Context, string) error          { return nil }
func (m *MockRoleRepo) AssignPermissions(context.Context, string, []string) error {
	return nil
}
func (m *MockRoleRepo) GetPermissionKeysByRoleID(_ context.Context, roleID string) ([]string, error) {
	return m.Permissions[roleID], nil
}
func (m *MockRoleRepo) GetPermissionIDsByKeys(context.Context, []string) ([]string, error) {
	return nil, nil
}

// MockUserRepo is a no-op UserRepository stub for auth unit tests.
type MockUserRepo struct{}

func (m *MockUserRepo) Create(context.Context, *usermodel.User) error { return nil }
func (m *MockUserRepo) GetByID(context.Context, string) (*usermodel.User, error) {
	return nil, nil
}
func (m *MockUserRepo) GetByEmail(context.Context, string) (*usermodel.User, error) {
	return nil, nil
}
func (m *MockUserRepo) List(context.Context, *query.Query) ([]usermodel.User, int64, error) {
	return nil, 0, nil
}
func (m *MockUserRepo) ListAll(context.Context) ([]usermodel.User, error) { return nil, nil }
func (m *MockUserRepo) Update(context.Context, *usermodel.User) error     { return nil }
func (m *MockUserRepo) Delete(context.Context, string) error              { return nil }
func (m *MockUserRepo) CountSuperAdmins(context.Context) (int64, error)   { return 0, nil }

// MemoryUserRepo is an in-memory UserRepository for user service unit tests.
type MemoryUserRepo struct {
	Users map[string]*usermodel.User
}

func (m *MemoryUserRepo) Create(_ context.Context, user *usermodel.User) error {
	if user.ID == "" {
		user.ID = "generated-id"
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MemoryUserRepo) GetByID(_ context.Context, id string) (*usermodel.User, error) {
	user, ok := m.Users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *MemoryUserRepo) GetByEmail(_ context.Context, email string) (*usermodel.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MemoryUserRepo) List(_ context.Context, q *query.Query) ([]usermodel.User, int64, error) {
	items := make([]usermodel.User, 0, len(m.Users))
	for _, user := range m.Users {
		items = append(items, *user)
	}
	total := int64(len(items))
	if q == nil {
		return items, total, nil
	}
	if q.Offset >= len(items) {
		return []usermodel.User{}, total, nil
	}
	end := q.Offset + q.Limit
	if end > len(items) {
		end = len(items)
	}
	return items[q.Offset:end], total, nil
}

func (m *MemoryUserRepo) ListAll(_ context.Context) ([]usermodel.User, error) {
	items := make([]usermodel.User, 0, len(m.Users))
	for _, user := range m.Users {
		items = append(items, *user)
	}
	return items, nil
}

func (m *MemoryUserRepo) Update(_ context.Context, user *usermodel.User) error {
	m.Users[user.ID] = user
	return nil
}

func (m *MemoryUserRepo) Delete(_ context.Context, id string) error {
	delete(m.Users, id)
	return nil
}

func (m *MemoryUserRepo) CountSuperAdmins(_ context.Context) (int64, error) {
	var n int64
	for _, user := range m.Users {
		if user != nil && user.IsSuperAdmin && user.DeletedAt == nil {
			n++
		}
	}
	return n, nil
}
