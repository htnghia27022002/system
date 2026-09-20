package maps

// SearchQuery filters GET /admin/maps/search.
type SearchQuery struct {
	Q           string `form:"q"`
	CountryCode string `form:"countryCode"`
}

// SearchHit is one saved Place, saved Location, or geocode result.
// Kind is place | location | geocode. No vendor names in the payload.
type SearchHit struct {
	ID         string  `json:"id"`
	Kind       string  `json:"kind"`
	Title      string  `json:"title"`
	Subtitle   string  `json:"subtitle"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	PlaceID    string  `json:"placeId,omitempty"`
	LocationID string  `json:"locationId,omitempty"`
}

// SearchResponse is the combined map search result.
type SearchResponse struct {
	Items []SearchHit `json:"items"`
}
