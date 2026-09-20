package utils

import "time"

const DateLayout = "2006-01-02"

// FormatRFC3339 returns t in UTC RFC3339.
func FormatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// FormatRFC3339Ptr returns UTC RFC3339 or nil.
func FormatRFC3339Ptr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := FormatRFC3339(*t)
	return &s
}

// FormatDate returns UTC YYYY-MM-DD or nil.
func FormatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(DateLayout)
	return &s
}
