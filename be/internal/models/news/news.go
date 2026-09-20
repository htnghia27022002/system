package news

import "time"

// News is one sourced article about a Place. Unique on original_url.
type News struct {
	ID          string         `json:"id" db:"id"`
	PlaceID     string         `json:"placeId" db:"place_id"`
	CategoryID  string         `json:"categoryId" db:"category_id"`
	SourceID    *string        `json:"sourceId" db:"source_id"`
	SourceName  string         `json:"sourceName" db:"source_name"`
	OriginalURL string         `json:"originalUrl" db:"original_url"`
	Title       string         `json:"title" db:"title"`
	Details     map[string]any `json:"details" db:"details"`
	CreatedBy   *string        `json:"createdBy" db:"created_by"`
	UpdatedBy   *string        `json:"updatedBy" db:"updated_by"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
}

func (News) TableName() string {
	return "news"
}
