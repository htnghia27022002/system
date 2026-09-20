package geo

import "time"

// Division is one node in a country's administrative tree.
type Division struct {
	ID        string    `json:"id" db:"id"`
	CountryID string    `json:"countryId" db:"country_id"`
	ParentID  *string   `json:"parentId" db:"parent_id"`
	Level     int16     `json:"level" db:"level"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	NameEn    string    `json:"nameEn" db:"name_en"`
	FullName  string    `json:"fullName" db:"full_name"`
	Type      string    `json:"type" db:"type"`
	Path      string    `json:"path" db:"path"`
	Lat       *float64  `json:"lat" db:"lat"`
	Lng       *float64  `json:"lng" db:"lng"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedBy *string   `json:"createdBy" db:"created_by"`
	UpdatedBy *string   `json:"updatedBy" db:"updated_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (Division) TableName() string {
	return "administrative_divisions"
}
