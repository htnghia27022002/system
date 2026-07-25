package migrations

import "embed"

// FS holds SQL migration files for golang-migrate (iofs).
// Embedded so a single deploy binary (e.g. Render) can migrate without a migrations/ directory on disk.
//
//go:embed *.sql
var FS embed.FS
