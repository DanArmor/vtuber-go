package embed_migrations

import "embed"

//go:embed ent/migrate/migrations
var Migrations embed.FS
