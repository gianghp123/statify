package migrations

import "embed"

// This works because the .sql files are in the SAME folder
//
//go:embed *.sql
var FS embed.FS
