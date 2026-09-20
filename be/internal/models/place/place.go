package place

import (
	"strings"
	"time"
)

const (
	StatusPending int16 = 1
	StatusActive  int16 = 2
	StatusHidden  int16 = 3
)

// Place is a reusable point of interest at a Location.
type Place struct {
	ID         string         `json:"id" db:"id"`
	LocationID string         `json:"locationId" db:"location_id"`
	CategoryID string         `json:"categoryId" db:"category_id"`
	PlaceKey   string         `json:"placeKey" db:"place_key"`
	Name       string         `json:"name" db:"name"`
	Unit       string         `json:"unit" db:"unit"`
	Status     int16          `json:"status" db:"status"`
	Lat        *float64       `json:"lat" db:"lat"`
	Lng        *float64       `json:"lng" db:"lng"`
	Details    map[string]any `json:"details" db:"details"`
	CreatedBy  *string        `json:"createdBy" db:"created_by"`
	UpdatedBy  *string        `json:"updatedBy" db:"updated_by"`
	CreatedAt  time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time      `json:"updatedAt" db:"updated_at"`
}

func (Place) TableName() string {
	return "places"
}

func StatusToJSON(status int16) string {
	switch status {
	case StatusActive:
		return "active"
	case StatusHidden:
		return "hidden"
	default:
		return "pending"
	}
}

func StatusFromJSON(status string) (int16, bool) {
	switch status {
	case "pending":
		return StatusPending, true
	case "active":
		return StatusActive, true
	case "hidden":
		return StatusHidden, true
	default:
		return 0, false
	}
}

// DerivePlaceKey returns a stable key for (location, place).
func DerivePlaceKey(name, unit, explicit string) string {
	if key := strings.ToLower(strings.TrimSpace(explicit)); key != "" {
		return key
	}
	return strings.ToLower(strings.TrimSpace(name)) + "|" + strings.ToLower(strings.TrimSpace(unit))
}

// ValidCoords reports whether both coordinates are present and in WGS84 range.
func ValidCoords(lat, lng *float64) bool {
	if lat == nil || lng == nil {
		return false
	}
	return *lat >= -90 && *lat <= 90 && *lng >= -180 && *lng <= 180
}
