package user

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBlocked  Status = "blocked"
)

type User struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;type:varchar(255)"`
	FullName     string         `json:"fullName" gorm:"column:full_name;type:varchar(255);not null"`
	Phone        string         `json:"phone" gorm:"type:varchar(50)"`
	AvatarURL    string         `json:"avatarUrl" gorm:"column:avatar_url;type:text"`
	RoleID       string         `json:"roleId" gorm:"column:role_id;type:uuid;index"`
	Status       Status         `json:"status" gorm:"type:varchar(20);not null;default:active;index"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
