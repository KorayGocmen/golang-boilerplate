package migrations

import "embed"

// Embed migrations into the binary
// to ship a single binary to production.
var (
	//go:embed sqls/*.sql
	FS   embed.FS
	Path = "sqls"
)
