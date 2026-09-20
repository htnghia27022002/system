package geo

import "time"

// Country is a shared ISO country catalog row.
type Country struct {
	ID        string    `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Code3     string    `json:"code3" db:"code3"`
	Name      string    `json:"name" db:"name"`
	NameLocal string    `json:"nameLocal" db:"name_local"`
	PhoneCode string    `json:"phoneCode" db:"phone_code"`
	Currency  string    `json:"currency" db:"currency"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedBy *string   `json:"createdBy" db:"created_by"`
	UpdatedBy *string   `json:"updatedBy" db:"updated_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (Country) TableName() string {
	return "countries"
}
