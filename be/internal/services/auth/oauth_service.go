package auth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	apperrors "be/internal/common/errors"
	"be/internal/config"
	authdto "be/internal/dto/auth"
	authmodel "be/internal/models/auth"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
	oauthprovider "be/internal/services/auth/oauth"
)

type OAuthService struct {
	cfg      config.Config
	authRepo interfaces.AuthRepository
	roleRepo interfaces.RoleRepository
	authSvc  *Service
	registry *oauthprovider.Registry
}

func NewOAuthService(cfg config.Config, authRepo interfaces.AuthRepository, roleRepo interfaces.RoleRepository, authSvc *Service) *OAuthService {
	return &OAuthService{
		cfg:      cfg,
		authRepo: authRepo,
		roleRepo: roleRepo,
		authSvc:  authSvc,
		registry: oauthprovider.NewRegistry(cfg),
	}
}

func (s *OAuthService) isProviderAllowed(provider string) bool {
	for _, item := range s.cfg.OAuthAllowedProviders {
		if item == provider {
			return true
		}
	}
	return false
}

func (s *OAuthService) provider(providerID string) (oauthprovider.Provider, error) {
	if !s.isProviderAllowed(providerID) {
		return nil, fmt.Errorf("%w: provider not allowed", apperrors.ErrBadRequest)
	}
	return s.registry.Get(providerID)
}

// IsProviderConfigured reports whether OAuth can be used for the provider (credentials present).
func (s *OAuthService) IsProviderConfigured(provider string) bool {
	p, err := s.provider(provider)
	if err != nil {
		return false
	}
	return p.Configured()
}

// ListConfiguredProviders returns provider ids that are allowed and configured.
func (s *OAuthService) ListConfiguredProviders() []string {
	out := make([]string, 0, len(s.cfg.OAuthAllowedProviders))
	for _, providerID := range s.cfg.OAuthAllowedProviders {
		if s.IsProviderConfigured(providerID) {
			out = append(out, providerID)
		}
	}
	return out
}

func (s *OAuthService) StartURL(providerID, redirectURI string) (string, error) {
	p, err := s.provider(providerID)
	if err != nil {
		return "", err
	}
	cfg, err := p.OAuthConfig(redirectURI)
	if err != nil {
		return "", err
	}
	return cfg.AuthCodeURL(providerID+"-state", oauth2.AccessTypeOffline), nil
}

func (s *OAuthService) Callback(ctx context.Context, providerID, code, redirectURI string) (*authdto.AuthResponse, error) {
	p, err := s.provider(providerID)
	if err != nil {
		return nil, err
	}
	cfg, err := p.OAuthConfig(redirectURI)
	if err != nil {
		return nil, err
	}

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: oauth code exchange failed", apperrors.ErrUnauthorized)
	}

	client := cfg.Client(ctx, token)
	profile, err := p.FetchProfile(ctx, client)
	if err != nil {
		return nil, err
	}

	user, err := s.resolveOAuthUser(ctx, providerID, profile, token)
	if err != nil {
		return nil, err
	}
	if user.Status != usermodel.StatusActive {
		return nil, fmt.Errorf("%w: account is not active", apperrors.ErrForbidden)
	}
	return s.authSvc.issueTokenPair(ctx, user)
}

func (s *OAuthService) resolveOAuthUser(ctx context.Context, provider string, profile *oauthprovider.Profile, token *oauth2.Token) (*usermodel.User, error) {
	account, err := s.authRepo.FindOAuthAccount(ctx, provider, profile.ID)
	if err != nil {
		return nil, err
	}

	var user *usermodel.User
	if account != nil {
		user, err = s.authRepo.FindUserByID(ctx, account.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, apperrors.ErrNotFound
		}
		s.syncOAuthAccount(ctx, account, profile, token)
		return user, nil
	}

	user, err = s.authRepo.FindUserByEmail(ctx, profile.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		defaultRole, err := s.roleRepo.GetBySlug(ctx, "user")
		if err != nil {
			return nil, err
		}
		if defaultRole == nil {
			return nil, fmt.Errorf("default member role is not configured")
		}
		user = &usermodel.User{
			Email:    profile.Email,
			FullName: profile.Name,
			RoleID:   defaultRole.ID,
			Status:   usermodel.StatusActive,
		}
		if err := s.authRepo.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	oauthAccount := &authmodel.OAuthAccount{
		UserID:         user.ID,
		Provider:       provider,
		ProviderUserID: profile.ID,
		ProviderEmail:  profile.Email,
		AccessToken:    token.AccessToken,
		TokenType:      token.TokenType,
	}
	if token.RefreshToken != "" {
		oauthAccount.RefreshToken = token.RefreshToken
	}
	if !token.Expiry.IsZero() {
		expiry := token.Expiry
		oauthAccount.ExpiresAt = &expiry
	}
	if err := s.authRepo.CreateOAuthAccount(ctx, oauthAccount); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *OAuthService) syncOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount, profile *oauthprovider.Profile, token *oauth2.Token) {
	account.ProviderEmail = profile.Email
	account.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		account.RefreshToken = token.RefreshToken
	}
	if !token.Expiry.IsZero() {
		expiry := token.Expiry
		account.ExpiresAt = &expiry
	}
	_ = s.authRepo.UpdateOAuthAccount(ctx, account)
}

// ValidateRedirectURI re-exports shared redirect validation for handlers.
func ValidateRedirectURI(raw string) error {
	return oauthprovider.ValidateRedirectURI(raw)
}
