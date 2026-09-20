package utils

import "strings"

// StringPtr returns a trimmed pointer, or nil when empty.
func StringPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
