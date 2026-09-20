package maps

// SourceRecord is the CRUD payload for a Data Source.
type SourceRecord struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Enabled      bool           `json:"enabled"`
	HTTPMethod   string         `json:"httpMethod"`
	URL          string         `json:"url"`
	Headers      map[string]any `json:"headers"`
	QueryParams  map[string]any `json:"queryParams"`
	Body         string         `json:"body"`
	FieldMapping map[string]any `json:"fieldMapping"`
	CreatedBy    *string        `json:"createdBy"`
	UpdatedBy    *string        `json:"updatedBy"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
}

// SourceListQuery paginates GET /sources.
type SourceListQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// SourceListResponse is the paginated Sources list.
type SourceListResponse struct {
	Items   []SourceRecord `json:"items"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
	Total   int64          `json:"total"`
	HasMore bool           `json:"hasMore"`
}

// CreateSourceRequest is POST /sources.
type CreateSourceRequest struct {
	Name         string         `json:"name"`
	Enabled      *bool          `json:"enabled"`
	HTTPMethod   string         `json:"httpMethod"`
	URL          string         `json:"url"`
	Headers      map[string]any `json:"headers"`
	QueryParams  map[string]any `json:"queryParams"`
	Body         string         `json:"body"`
	FieldMapping map[string]any `json:"fieldMapping"`
}

// ProbeSourceRequest is POST /sources/probe — fetch an unsaved URL to preview JSON.
type ProbeSourceRequest struct {
	HTTPMethod  string         `json:"httpMethod"`
	URL         string         `json:"url"`
	Headers     map[string]any `json:"headers"`
	QueryParams map[string]any `json:"queryParams"`
	Body        string         `json:"body"`
}

// ProbeSourceResponse is the JSON body returned by a successful probe.
type ProbeSourceResponse struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}

// PatchSourceRequest is PATCH /sources/:id (all fields optional).
type PatchSourceRequest struct {
	Name         *string        `json:"name"`
	Enabled      *bool          `json:"enabled"`
	HTTPMethod   *string        `json:"httpMethod"`
	URL          *string        `json:"url"`
	Headers      map[string]any `json:"headers"`
	QueryParams  map[string]any `json:"queryParams"`
	Body         *string        `json:"body"`
	FieldMapping map[string]any `json:"fieldMapping"`
}
