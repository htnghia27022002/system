package auth

import "time"

type RefreshToken struct {
	ID         string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     string     `json:"userId" gorm:"column:user_id;type:uuid;not null;index"`
	TokenHash  string     `json:"-" gorm:"column:token_hash;type:varchar(255);not null;uniqueIndex"`
	ExpiresAt  time.Time  `json:"expiresAt" gorm:"column:expires_at;not null;index"`
	RevokedAt  *time.Time `json:"revokedAt,omitempty" gorm:"column:revoked_at"`
	CreatedAt  time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	IPAddress  string     `json:"-" gorm:"column:ip_address;type:varchar(64)"`
	UserAgent  string     `json:"-" gorm:"column:user_agent;type:text"`
	DeviceID   string     `json:"-" gorm:"column:device_id;type:varchar(255)"`
	RawToken   string     `json:"-" gorm:"-"`
}

type OAuthAccount struct {
	ID             string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         string     `json:"userId" gorm:"column:user_id;type:uuid;not null;index"`
	Provider       string     `json:"provider" gorm:"type:varchar(50);not null;uniqueIndex:idx_oauth_provider_user"`
	ProviderUserID string     `json:"providerUserId" gorm:"column:provider_user_id;type:varchar(255);not null;uniqueIndex:idx_oauth_provider_user"`
	ProviderEmail  string     `json:"providerEmail,omitempty" gorm:"column:provider_email;type:varchar(255)"`
	AccessToken    string     `json:"-" gorm:"column:access_token;type:text"`
	RefreshToken   string     `json:"-" gorm:"column:refresh_token;type:text"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty" gorm:"column:expires_at"`
	Scope          string     `json:"scope,omitempty" gorm:"type:text"`
	TokenType      string     `json:"tokenType,omitempty" gorm:"column:token_type;type:varchar(50)"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}
