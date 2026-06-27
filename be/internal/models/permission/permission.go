package permission

import "time"

type Permission struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Key         string    `json:"key" gorm:"column:key;type:varchar(100);not null;uniqueIndex"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Group       string    `json:"group" gorm:"column:group;type:varchar(50);not null;default:''"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
