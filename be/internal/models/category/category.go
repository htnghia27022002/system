package category

import "time"

const (
	KeyRoomRental    = "room_rental"
	KeyRestaurant    = "restaurant"
	KeyHotel         = "hotel"
	KeyEatery        = "eatery"
	KeyUncategorized = "uncategorized"
)

// Category is a shared danh muc catalog row. Unique on key.
type Category struct {
	ID        string    `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Name      string    `json:"name" db:"name"`
	NameLocal string    `json:"nameLocal" db:"name_local"`
	SortOrder int       `json:"sortOrder" db:"sort_order"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedBy *string   `json:"createdBy" db:"created_by"`
	UpdatedBy *string   `json:"updatedBy" db:"updated_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

// ToJSONKey maps a stored snake_case key to the API camelCase token.
func ToJSONKey(key string) string {
	if key == KeyRoomRental {
		return "roomRental"
	}
	return key
}

// FromJSONKey maps an API category token to the stored key.
func FromJSONKey(key string) string {
	if key == "roomRental" {
		return KeyRoomRental
	}
	return key
}

// IsP1FilterKey reports whether key is a valid GET /places category filter.
func IsP1FilterKey(key string) bool {
	switch FromJSONKey(key) {
	case KeyRoomRental, KeyRestaurant, KeyHotel, KeyEatery:
		return true
	default:
		return false
	}
}

// IsCatalogKey reports whether key is a stored or JSON catalog key.
func IsCatalogKey(key string) bool {
	switch FromJSONKey(key) {
	case KeyRoomRental, KeyRestaurant, KeyHotel, KeyEatery, KeyUncategorized:
		return true
	default:
		return false
	}
}
