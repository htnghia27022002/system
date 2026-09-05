package postgres

import "strings"

// QuoteIdent quotes a PostgreSQL identifier. Names that are already quoted are left unchanged.
func QuoteIdent(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || strings.HasPrefix(name, `"`) {
		return name
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
