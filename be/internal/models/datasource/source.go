package datasource

import (
	"strings"
	"time"
)

const (
	HTTPGet   int16 = 1
	HTTPPost  int16 = 2
	HTTPPut   int16 = 3
	HTTPPatch int16 = 4
)

// Source is an operator-managed HTTP ingest endpoint.
type Source struct {
	ID           string         `json:"id" db:"id"`
	Name         string         `json:"name" db:"name"`
	Enabled      bool           `json:"enabled" db:"enabled"`
	HTTPMethod   int16          `json:"httpMethod" db:"http_method"`
	URL          string         `json:"url" db:"url"`
	Headers      map[string]any `json:"headers" db:"headers"`
	QueryParams  map[string]any `json:"queryParams" db:"query_params"`
	Body         string         `json:"body" db:"body"`
	FieldMapping map[string]any `json:"fieldMapping" db:"field_mapping"`
	CreatedBy    *string        `json:"createdBy" db:"created_by"`
	UpdatedBy    *string        `json:"updatedBy" db:"updated_by"`
	CreatedAt    time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time      `json:"updatedAt" db:"updated_at"`
}

func (Source) TableName() string {
	return "data_sources"
}

func MethodToJSON(method int16) string {
	switch method {
	case HTTPPost:
		return "POST"
	case HTTPPut:
		return "PUT"
	case HTTPPatch:
		return "PATCH"
	default:
		return "GET"
	}
}

func MethodFromJSON(method string) (int16, bool) {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case "GET":
		return HTTPGet, true
	case "POST":
		return HTTPPost, true
	case "PUT":
		return HTTPPut, true
	case "PATCH":
		return HTTPPatch, true
	default:
		return 0, false
	}
}
