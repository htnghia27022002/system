// Package migrations holds the webhooks module schema (own version table).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
