// Package db exposes the embedded migration files so they can be used by
// the migration runner and tests without duplicating the embed directive.
package db

import "embed"

// Migrations holds all goose SQL migration files embedded at compile time.
//
//go:embed versions/*.sql
var Migrations embed.FS
