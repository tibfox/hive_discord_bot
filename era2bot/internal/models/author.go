package models

import "time"

//go:generate stringer -type=Status

type Author struct {
	Name            string
	BlacklistedOn   time.Time
	BlacklistedBy   string
	BlacklistReason string
	// TODO: add more here like amount of curations past 7 days and such
}
