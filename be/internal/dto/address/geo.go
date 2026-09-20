package address

// CountryItem is one active country in GET /address/countries.
type CountryItem struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Code3     string `json:"code3"`
	Name      string `json:"name"`
	NameLocal string `json:"nameLocal"`
}

// CountryListResponse is the countries catalog payload.
type CountryListResponse struct {
	Items []CountryItem `json:"items"`
}

// AdminDivision is one administrative unit in GET /address/divisions.
type AdminDivision struct {
	ID          string   `json:"id"`
	CountryCode string   `json:"countryCode"`
	ParentID    *string  `json:"parentId"`
	Level       int16    `json:"level"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	NameEn      string   `json:"nameEn"`
	FullName    string   `json:"fullName"`
	Type        string   `json:"type"`
	Path        string   `json:"path"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
}

// DivisionListQuery filters GET /address/divisions.
type DivisionListQuery struct {
	CountryCode string `form:"countryCode"`
	ParentID    string `form:"parentId"`
	Q           string `form:"q"`
}

// DivisionListResponse is the divisions catalog payload.
type DivisionListResponse struct {
	Items []AdminDivision `json:"items"`
}
