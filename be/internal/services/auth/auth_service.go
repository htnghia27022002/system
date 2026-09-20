package auth

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	apperrors "be/common/errors"
	jwtmanager "be/common/jwt"
	"be/common/utils"
	authdto "be/internal/dto/auth"
	userdto "be/internal/dto/user"
	authmodel "be/internal/models/auth"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
	"be/internal/services/media"
	"be/pkg/hash"
)

type Service struct {
	authRepo   interfaces.AuthRepository
	userRepo   interfaces.UserRepository
	roleRepo   interfaces.RoleRepository
	jwt        *jwtmanager.Manager
	refreshTTL time.Duration
	media      media.AvatarStorage
}

func NewService(
	authRepo interfaces.AuthRepository,
	userRepo interfaces.UserRepository,
	roleRepo interfaces.RoleRepository,
	jwt *jwtmanager.Manager,
	refreshTTL time.Duration,
	mediaSvc media.AvatarStorage,
) *Service {
	return &Service{
		authRepo:   authRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwt:        jwt,
		refreshTTL: refreshTTL,
		media:      mediaSvc,
	}
}

func (s *Service) ResolveAuthUser(ctx context.Context, userID string) (*authdto.AuthUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}
	return s.buildAuthUser(ctx, user)
}

func (s *Service) buildAuthUser(ctx context.Context, user *usermodel.User) (*authdto.AuthUserResponse, error) {
	roleSlug := "user"
	if user.RoleID != "" {
		role, err := s.roleRepo.GetByID(ctx, user.RoleID)
		if err != nil {
			return nil, err
		}
		if role != nil {
			roleSlug = role.Slug
		}
	}

	permissions, err := s.roleRepo.GetPermissionKeysByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}
	if permissions == nil {
		permissions = []string{}
	}

	return &authdto.AuthUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.FullName,
		Role:        roleSlug,
		RoleID:      user.RoleID,
		Permissions: permissions,
		SuperAdmin:  user.IsSuperAdmin,
		Phone:       user.Phone,
		AvatarURL:   user.AvatarURL,
		General:     user.General,
		Birthday:    utils.FormatDate(user.Birthday),
		Address:     user.Address,
		SocialLinks: userdto.SocialLinksToDTO(user.SocialLinks),
		HasPassword: user.PasswordHash != "",
	}, nil
}

func (s *Service) issueTokenPair(ctx context.Context, user *usermodel.User, ip, userAgent string) (*authdto.AuthResponse, error) {
	authUser, err := s.buildAuthUser(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.SignAccessToken(
		user.ID,
		user.Email,
		user.FullName,
		authUser.Role,
		user.RoleID,
		authUser.Permissions,
		authUser.SuperAdmin,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	rawRefresh := uuid.NewString()
	refresh := &authmodel.RefreshToken{
		UserID:     user.ID,
		RawToken:   rawRefresh,
		TokenHash:  hash.SHA256Hex(rawRefresh),
		ExpiresAt:  now.Add(s.refreshTTL),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	applySessionMeta(refresh, ip, userAgent, now)
	if err := s.authRepo.CreateRefreshToken(ctx, refresh); err != nil {
		return nil, err
	}

	return &authdto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		SessionID:    refresh.ID,
		User:         *authUser,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password, ip, userAgent string) (*authdto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("%w: invalid credentials", apperrors.ErrUnauthorized)
	}
	if user.PasswordHash == "" {
		return nil, fmt.Errorf("%w: invalid credentials", apperrors.ErrUnauthorized)
	}
	if err := hash.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", apperrors.ErrUnauthorized)
	}
	if user.Status != usermodel.StatusActive {
		return nil, fmt.Errorf("%w: account is not active", apperrors.ErrForbidden)
	}
	return s.issueTokenPair(ctx, user, ip, userAgent)
}

func (s *Service) Register(ctx context.Context, name, email, password, ip, userAgent string) (*authdto.AuthResponse, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: email already registered", apperrors.ErrConflict)
	}

	defaultRole, err := s.roleRepo.GetBySlug(ctx, "user")
	if err != nil {
		return nil, err
	}
	if defaultRole == nil {
		return nil, errors.New("default member role is not configured")
	}

	hashed, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &usermodel.User{
		Email:        email,
		PasswordHash: hashed,
		FullName:     name,
		RoleID:       defaultRole.ID,
		Status:       usermodel.StatusActive,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return s.issueTokenPair(ctx, user, ip, userAgent)
}

func (s *Service) Refresh(ctx context.Context, refreshToken, ip, userAgent string) (*authdto.TokenPairResponse, error) {
	stored, err := s.authRepo.FindRefreshTokenByHash(ctx, hash.SHA256Hex(refreshToken))
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, fmt.Errorf("%w: invalid refresh token", apperrors.ErrUnauthorized)
	}

	user, err := s.userRepo.GetByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("%w: invalid refresh token", apperrors.ErrUnauthorized)
	}

	if err := s.authRepo.RevokeRefreshToken(ctx, stored.TokenHash, time.Now()); err != nil {
		return nil, err
	}

	authUser, err := s.buildAuthUser(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.SignAccessToken(
		user.ID,
		user.Email,
		user.FullName,
		authUser.Role,
		user.RoleID,
		authUser.Permissions,
		authUser.SuperAdmin,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	rawRefresh := uuid.NewString()
	rotated := &authmodel.RefreshToken{
		UserID:     stored.UserID,
		RawToken:   rawRefresh,
		TokenHash:  hash.SHA256Hex(rawRefresh),
		ExpiresAt:  now.Add(s.refreshTTL),
		CreatedAt:  sessionStartedAt(stored.CreatedAt, now),
		UpdatedAt:  now,
	}
	applySessionMeta(rotated, ip, userAgent, now)
	if err := s.authRepo.CreateRefreshToken(ctx, rotated); err != nil {
		return nil, err
	}

	return &authdto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		SessionID:    rotated.ID,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	stored, err := s.authRepo.FindRefreshTokenByHash(ctx, hash.SHA256Hex(refreshToken))
	if err != nil {
		return err
	}
	if stored == nil {
		return nil
	}
	return s.authRepo.RevokeRefreshToken(ctx, stored.TokenHash, time.Now())
}

// ListSessions returns the JWT owner's active sessions. currentSessionID is the
// refresh_tokens.id from X-Session-Id (not the refresh secret).
func (s *Service) ListSessions(ctx context.Context, userID, currentSessionID string) (*authdto.SessionListResponse, error) {
	rows, err := s.authRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	currentSessionID = strings.TrimSpace(currentSessionID)
	items := make([]authdto.SessionListItem, 0, len(rows))
	for i := range rows {
		items = append(items, mapSessionListItem(&rows[i], currentSessionID))
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Current != items[j].Current {
			return items[i].Current
		}
		return items[i].LastUsedAt.After(items[j].LastUsedAt)
	})
	return &authdto.SessionListResponse{Items: items}, nil
}

// RevokeSession ends one owned active session. 400 if target is the current session; 404 if missing/not owned/inactive.
func (s *Service) RevokeSession(ctx context.Context, userID, targetID, currentSessionID string) error {
	targetID = strings.TrimSpace(targetID)
	currentSessionID = strings.TrimSpace(currentSessionID)
	if targetID == "" {
		return fmt.Errorf("%w: session id is required", apperrors.ErrBadRequest)
	}
	if currentSessionID != "" && targetID == currentSessionID {
		return fmt.Errorf("%w: cannot revoke the current session; use sign out", apperrors.ErrBadRequest)
	}
	n, err := s.authRepo.RevokeByID(ctx, userID, targetID)
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%w: session not found", apperrors.ErrNotFound)
	}
	return nil
}

// RevokeOtherSessions revokes every active session except sessionID. 400 if sessionID is missing or not an owned active session.
func (s *Service) RevokeOtherSessions(ctx context.Context, userID, sessionID string) (*authdto.SessionListResponse, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, fmt.Errorf("%w: sessionId is required", apperrors.ErrBadRequest)
	}
	rows, err := s.authRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	owned := false
	for i := range rows {
		if rows[i].ID == sessionID {
			owned = true
			break
		}
	}
	if !owned {
		return nil, fmt.Errorf("%w: session is not an active session of this user", apperrors.ErrBadRequest)
	}
	if _, err := s.authRepo.RevokeAllExcept(ctx, userID, sessionID); err != nil {
		return nil, err
	}
	return s.ListSessions(ctx, userID, sessionID)
}

func applySessionMeta(token *authmodel.RefreshToken, ip, userAgent string, lastUsed time.Time) {
	token.IPAddress = clipClientIP(ip)
	token.UserAgent = strings.TrimSpace(userAgent)
	token.LastUsedAt = lastUsed
}

func sessionStartedAt(stored time.Time, fallback time.Time) time.Time {
	if stored.IsZero() || stored.Unix() <= 0 {
		return fallback
	}
	return stored
}

func clipClientIP(ip string) string {
	ip = strings.TrimSpace(ip)
	if len(ip) > 64 {
		return ip[:64]
	}
	return ip
}

func mapSessionListItem(row *authmodel.RefreshToken, currentID string) authdto.SessionListItem {
	started := sessionStartedAt(row.CreatedAt, row.LastUsedAt)
	if started.IsZero() || started.Unix() <= 0 {
		started = row.ExpiresAt
	}
	return authdto.SessionListItem{
		ID:         row.ID,
		CreatedAt:  started.UTC(),
		ExpiresAt:  row.ExpiresAt,
		LastUsedAt: row.LastUsedAt,
		IPAddress:  nullableMeta(row.IPAddress),
		UserAgent:  nullableMeta(row.UserAgent),
		Current:    currentID != "" && row.ID == currentID,
	}
}

func nullableMeta(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func (s *Service) Me(ctx context.Context, userID string) (*authdto.AuthUserResponse, error) {
	return s.ResolveAuthUser(ctx, userID)
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, req authdto.UpdateProfileRequest) (*authdto.AuthUserResponse, error) {
	if err := userdto.ValidateName(req.Name); err != nil {
		return nil, err
	}
	normalized, err := userdto.ValidatePersonalFields(req.Phone, req.General, req.Birthday, req.Address, req.SocialLinks)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	user.FullName = strings.TrimSpace(req.Name)
	normalized.Apply(user)
	if user.SocialLinks == nil {
		user.SocialLinks = []usermodel.SocialLink{}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return s.buildAuthUser(ctx, user)
}

func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("%w: new password must be at least 8 characters", apperrors.ErrBadRequest)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrNotFound
	}
	if user.PasswordHash == "" {
		return fmt.Errorf("%w: password change is not available for this account", apperrors.ErrBadRequest)
	}
	if err := hash.ComparePassword(user.PasswordHash, currentPassword); err != nil {
		return fmt.Errorf("%w: current password is incorrect", apperrors.ErrUnauthorized)
	}

	hashed, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hashed
	return s.userRepo.Update(ctx, user)
}

func (s *Service) UploadAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (*authdto.AuthUserResponse, error) {
	if s.media == nil {
		return nil, fmt.Errorf("%w: media storage is not configured", apperrors.ErrBadRequest)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	publicPath, err := s.media.SaveAvatar(file, header)
	if err != nil {
		return nil, err
	}

	prev := user.AvatarURL
	user.AvatarURL = publicPath
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.media.DeleteByPublicPath(publicPath)
		return nil, err
	}
	if prev != "" && prev != publicPath {
		s.media.DeleteByPublicPath(prev)
	}
	return s.buildAuthUser(ctx, user)
}
