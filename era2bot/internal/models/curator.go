package models

//go:generate stringer -type=Status

type Curator struct {
	Name                     string
	DiscordUserId            string
	CountCurationsInInterval int
}
