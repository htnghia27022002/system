package maps

// LocationListQuery filters GET /locations.
type LocationListQuery struct {
	Q               string `form:"q"`
	CountryCode     string `form:"countryCode"`
	AdminDivisionID string `form:"adminDivisionId"`
	Page            int    `form:"page"`
	Limit           int    `form:"limit"`
}

// LocationRecord is one Location for operator CRUD.
type LocationRecord struct {
	LocationSummary
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// LocationListResponse is the paginated Location list.
type LocationListResponse struct {
	Items   []LocationRecord `json:"items"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
	Total   int64            `json:"total"`
	HasMore bool             `json:"hasMore"`
}

// CreateLocationRequest creates a Location site.
type CreateLocationRequest struct {
	Name            string   `json:"name"`
	LocationKey     *string  `json:"locationKey"`
	CountryCode     *string  `json:"countryCode"`
	AdminDivisionID *string  `json:"adminDivisionId"`
	Street          string   `json:"street"`
	PostalCode      string   `json:"postalCode"`
	Formatted       string   `json:"formatted"`
	Lat             *float64 `json:"lat"`
	Lng             *float64 `json:"lng"`
}

// PatchLocationRequest updates a Location site.
type PatchLocationRequest struct {
	Name            *string  `json:"name"`
	LocationKey     *string  `json:"locationKey"`
	CountryCode     *string  `json:"countryCode"`
	AdminDivisionID *string  `json:"adminDivisionId"`
	Street          *string  `json:"street"`
	PostalCode      *string  `json:"postalCode"`
	Formatted       *string  `json:"formatted"`
	Lat             *float64 `json:"lat"`
	Lng             *float64 `json:"lng"`
}
