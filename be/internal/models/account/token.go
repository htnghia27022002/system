package account

import "time"

// TokenPurpose scopes a single-use user token to one flow.
type TokenPurpose string

const (
	PurposePasswordReset TokenPurpose = "password_reset"
	PurposeEmailVerify   TokenPurpose = "email_verify"
	PurposeInvite        TokenPurpose = "invite"
)

// UserToken is a hashed, expiring, single-use token sent by email.
type UserToken struct {
	ID        string       `json:"id" db:"id"`
	UserID    string       `json:"userId" db:"user_id"`
	Purpose   TokenPurpose `json:"purpose" db:"purpose"`
	TokenHash string       `json:"-" db:"token_hash"`
	ExpiresAt time.Time    `json:"expiresAt" db:"expires_at"`
	UsedAt    *time.Time   `json:"usedAt,omitempty" db:"used_at"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
}
