package maps

// PlaceListQuery filters GET /places.
type PlaceListQuery struct {
	Category        string `form:"category"`
	Q               string `form:"q"`
	AdminDivisionID string `form:"adminDivisionId"`
	CountryCode     string `form:"countryCode"`
	LocationID      string `form:"locationId"`
	Status          string `form:"status"`
	Manage          bool   `form:"manage"`
	Page            int    `form:"page"`
	Limit           int    `form:"limit"`
}

// PlacePin is one pinnable Place on the map list.
type PlacePin struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	CategoryID   string   `json:"categoryId"`
	Category     string   `json:"category"`
	Status       string   `json:"status"`
	Lat          *float64 `json:"lat"`
	Lng          *float64 `json:"lng"`
	LocationID   string   `json:"locationId"`
	LocationName string   `json:"locationName"`
	Unit         string   `json:"unit,omitempty"`
	PlaceKey     string   `json:"placeKey,omitempty"`
}

// PlaceListResponse is the paginated pin list.
type PlaceListResponse struct {
	Items   []PlacePin `json:"items"`
	Page    int        `json:"page"`
	Limit   int        `json:"limit"`
	Total   int64      `json:"total"`
	HasMore bool       `json:"hasMore"`
}

// LocationSummary is the Location identity on Place detail.
type LocationSummary struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	LocationKey     *string  `json:"locationKey"`
	CountryCode     *string  `json:"countryCode"`
	AdminDivisionID *string  `json:"adminDivisionId"`
	AdminPath       *string  `json:"adminPath"`
	Street          string   `json:"street"`
	PostalCode      string   `json:"postalCode"`
	Formatted       string   `json:"formatted"`
	Lat             *float64 `json:"lat"`
	Lng             *float64 `json:"lng"`
}

// NewsItem is one article on Place detail.
type NewsItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	OriginalURL string  `json:"originalUrl"`
	CategoryID  string  `json:"categoryId"`
	Category    string  `json:"category"`
	SourceName  string  `json:"sourceName"`
	SourceID    *string `json:"sourceId"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// PlaceDetail is GET/PATCH /places/:id.
type PlaceDetail struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Unit       string          `json:"unit"`
	PlaceKey   string          `json:"placeKey"`
	CategoryID string          `json:"categoryId"`
	Category   string          `json:"category"`
	Status     string          `json:"status"`
	Lat        *float64        `json:"lat"`
	Lng        *float64        `json:"lng"`
	Details    map[string]any  `json:"details"`
	Location   LocationSummary `json:"location"`
	News       []NewsItem      `json:"news"`
	CreatedAt  string          `json:"createdAt"`
	UpdatedAt  string          `json:"updatedAt"`
}

// CreatePlaceRequest creates a Place at a Location.
type CreatePlaceRequest struct {
	LocationID string         `json:"locationId"`
	Name       string         `json:"name"`
	Category   string         `json:"category"`
	PlaceKey   string         `json:"placeKey"`
	Unit       string         `json:"unit"`
	Status     string         `json:"status"`
	Lat        *float64       `json:"lat"`
	Lng        *float64       `json:"lng"`
	Details    map[string]any `json:"details"`
}

// PatchPlaceRequest updates a Place (status and/or fields).
type PatchPlaceRequest struct {
	LocationID *string         `json:"locationId"`
	Name       *string         `json:"name"`
	Category   *string         `json:"category"`
	PlaceKey   *string         `json:"placeKey"`
	Unit       *string         `json:"unit"`
	Status     *string         `json:"status"`
	Lat        *float64        `json:"lat"`
	Lng        *float64        `json:"lng"`
	Details    *map[string]any `json:"details"`
}
