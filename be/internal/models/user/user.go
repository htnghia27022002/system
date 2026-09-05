package user

import "time"

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBlocked  Status = "blocked"
)

// SocialLink is one optional labeled URL on a user profile (max 5 per user).
type SocialLink struct {
	Label string `json:"label,omitempty"`
	URL   string `json:"url"`
}

type User struct {
	ID           string       `json:"id" db:"id"`
	Email        string       `json:"email" db:"email"`
	PasswordHash string       `json:"-" db:"password_hash"`
	FullName     string       `json:"fullName" db:"full_name"`
	Phone        string       `json:"phone" db:"phone"`
	AvatarURL    string       `json:"avatarUrl" db:"avatar_url"`
	General      string       `json:"general" db:"general"`
	Birthday     *time.Time   `json:"birthday,omitempty" db:"birthday"`
	Address      string       `json:"address" db:"address"`
	SocialLinks  []SocialLink `json:"socialLinks" db:"social_links"`
	RoleID       string       `json:"roleId" db:"role_id"`
	IsSuperAdmin bool         `json:"superAdmin" db:"is_super_admin"`
	Status       Status       `json:"status" db:"status"`
	CreatedAt    time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt    *time.Time   `json:"-" db:"deleted_at"`
}
