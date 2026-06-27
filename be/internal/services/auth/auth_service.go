package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	apperrors "be/internal/common/errors"
	"be/internal/common/hash"
	jwtmanager "be/internal/common/jwt"
	authdto "be/internal/dto/auth"
	authmodel "be/internal/models/auth"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
)

type Service struct {
	authRepo interfaces.AuthRepository
	userRepo interfaces.UserRepository
	roleRepo interfaces.RoleRepository
	jwt      *jwtmanager.Manager
	refreshTTL time.Duration
}

func NewService(
	authRepo interfaces.AuthRepository,
	userRepo interfaces.UserRepository,
	roleRepo interfaces.RoleRepository,
	jwt *jwtmanager.Manager,
	refreshTTL time.Duration,
) *Service {
	return &Service{
		authRepo:   authRepo,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwt:        jwt,
		refreshTTL: refreshTTL,
	}
}

func hashRefreshToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func (s *Service) ResolveAuthUser(ctx context.Context, userID string) (*authdto.AuthUserResponse, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
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
	}, nil
}

func (s *Service) issueTokenPair(ctx context.Context, user *usermodel.User) (*authdto.AuthResponse, error) {
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
	)
	if err != nil {
		return nil, err
	}

	rawRefresh := uuid.NewString()
	refresh := &authmodel.RefreshToken{
		UserID:    user.ID,
		RawToken:  rawRefresh,
		TokenHash: hashRefreshToken(rawRefresh),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, refresh); err != nil {
		return nil, err
	}

	return &authdto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		User:         *authUser,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*authdto.AuthResponse, error) {
	user, err := s.authRepo.FindUserByEmail(ctx, email)
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
	return s.issueTokenPair(ctx, user)
}

func (s *Service) Register(ctx context.Context, name, email, password string) (*authdto.AuthResponse, error) {
	existing, err := s.authRepo.FindUserByEmail(ctx, email)
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
	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return s.issueTokenPair(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*authdto.TokenPairResponse, error) {
	stored, err := s.authRepo.FindRefreshTokenByHash(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, fmt.Errorf("%w: invalid refresh token", apperrors.ErrUnauthorized)
	}

	user, err := s.authRepo.FindUserByID(ctx, stored.UserID)
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
	)
	if err != nil {
		return nil, err
	}

	rawRefresh := uuid.NewString()
	rotated := &authmodel.RefreshToken{
		UserID:    stored.UserID,
		RawToken:  rawRefresh,
		TokenHash: hashRefreshToken(rawRefresh),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.authRepo.CreateRefreshToken(ctx, rotated); err != nil {
		return nil, err
	}

	return &authdto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	stored, err := s.authRepo.FindRefreshTokenByHash(ctx, hashRefreshToken(refreshToken))
	if err != nil {
		return err
	}
	if stored == nil {
		return nil
	}
	return s.authRepo.RevokeRefreshToken(ctx, stored.TokenHash, time.Now())
}

func (s *Service) Me(ctx context.Context, userID string) (*authdto.AuthUserResponse, error) {
	return s.ResolveAuthUser(ctx, userID)
}
