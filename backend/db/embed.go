// Package db exposes the embedded migration files so they can be used by
// both the database package and tests without duplicating the embed directive.
package db

import "embed"

// Migrations holds all goose SQL migration files embedded at compile time.
//
//go:embed migrations/*.sql
var Migrations embed.FS
