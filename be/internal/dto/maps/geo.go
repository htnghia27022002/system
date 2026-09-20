package maps

const (
	DefaultListLimit = 50
	MaxListLimit     = 200
)

// CategoryItem is one active category in GET /categories.
type CategoryItem struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	NameLocal string `json:"nameLocal"`
	SortOrder int    `json:"sortOrder"`
}

// CategoryListResponse is the categories catalog payload.
type CategoryListResponse struct {
	Items []CategoryItem `json:"items"`
}

// NormalizeList clamps page/limit to contract defaults (page=1, limit=50, max 200).
func NormalizeList(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}
	return page, limit
}
