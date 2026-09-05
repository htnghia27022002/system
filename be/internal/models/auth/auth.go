package auth

import "time"

type RefreshToken struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"userId" db:"user_id"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expiresAt" db:"expires_at"`
	RevokedAt *time.Time `json:"revokedAt,omitempty" db:"revoked_at"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
	LastUsedAt time.Time  `json:"lastUsedAt" db:"last_used_at"`
	IPAddress  string     `json:"-" db:"ip_address"`
	UserAgent  string     `json:"-" db:"user_agent"`
	DeviceID   string     `json:"-" db:"device_id"`
	RawToken   string     `json:"-" db:"-"`
}

type OAuthAccount struct {
	ID             string     `json:"id" db:"id"`
	UserID         string     `json:"userId" db:"user_id"`
	Provider       string     `json:"provider" db:"provider"`
	ProviderUserID string     `json:"providerUserId" db:"provider_user_id"`
	ProviderEmail  string     `json:"providerEmail,omitempty" db:"provider_email"`
	AccessToken    string     `json:"-" db:"access_token"`
	RefreshToken   string     `json:"-" db:"refresh_token"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	Scope          string     `json:"scope,omitempty" db:"scope"`
	TokenType      string     `json:"tokenType,omitempty" db:"token_type"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}
