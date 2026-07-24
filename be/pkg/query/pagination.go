package query

import "math"

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// NormalizePage clamps page and pageSize to supported defaults and limits.
func NormalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}

// Offset returns SQL offset for the given page and page size.
func Offset(page, pageSize int) int {
	page, pageSize = NormalizePage(page, pageSize)
	return (page - 1) * pageSize
}

// TotalPages returns total pages for a result count.
func TotalPages(total int64, pageSize int) int {
	if total == 0 || pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
