package auth

import (
	"time"

	userdto "be/internal/dto/user"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type UpdateProfileRequest struct {
	Name        string                   `json:"name" binding:"required,min=2"`
	Phone       *string                  `json:"phone"`
	General     *string                  `json:"general"`
	Birthday    *string                  `json:"birthday"`
	Address     *string                  `json:"address"`
	SocialLinks *[]userdto.SocialLinkDTO `json:"socialLinks"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

type AuthUserResponse struct {
	ID          string                  `json:"id"`
	Email       string                  `json:"email"`
	Name        string                  `json:"name"`
	Role        string                  `json:"role"`
	RoleID      string                  `json:"roleId"`
	Permissions []string                `json:"permissions"`
	SuperAdmin  bool                    `json:"superAdmin"`
	Phone       string                  `json:"phone"`
	AvatarURL   string                  `json:"avatarUrl"`
	General     string                  `json:"general"`
	Birthday    *string                 `json:"birthday"`
	Address     string                  `json:"address"`
	SocialLinks []userdto.SocialLinkDTO `json:"socialLinks"`
	HasPassword bool                    `json:"hasPassword"`
}

type AuthResponse struct {
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	SessionID    string           `json:"sessionId"`
	User         AuthUserResponse `json:"user"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	SessionID    string `json:"sessionId"`
}

// SessionListItem is one active refresh-token session for the JWT owner.
type SessionListItem struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	IPAddress  *string   `json:"ipAddress"`
	UserAgent  *string   `json:"userAgent"`
	Current    bool      `json:"current"`
}

type SessionListResponse struct {
	Items []SessionListItem `json:"items"`
}

type RevokeOthersRequest struct {
	SessionID string `json:"sessionId"`
}

type OAuthCallbackRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirectUri" binding:"required,url"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
