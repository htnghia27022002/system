package location

import "time"

// Location is a reusable geographic site (building, campus, depot).
type Location struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	LocationKey     *string   `json:"locationKey" db:"location_key"`
	CountryID       *string   `json:"countryId" db:"country_id"`
	AdminDivisionID *string   `json:"adminDivisionId" db:"admin_division_id"`
	Street          string    `json:"street" db:"street"`
	PostalCode      string    `json:"postalCode" db:"postal_code"`
	Formatted       string    `json:"formatted" db:"formatted"`
	Lat             *float64  `json:"lat" db:"lat"`
	Lng             *float64  `json:"lng" db:"lng"`
	CreatedBy       *string   `json:"createdBy" db:"created_by"`
	UpdatedBy       *string   `json:"updatedBy" db:"updated_by"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

func (Location) TableName() string {
	return "locations"
}
