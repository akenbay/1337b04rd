package domain

import "time"

type User struct {
	SessionID     string  // UUID
	AvatarURL     string  // From Rick & Morty API
	CharacterName string  // From API
	CustomName    *string // User can override (nullable)
	CreatedAt     time.Time
	ExpiresAt     time.Time // Session expiry (1 week)
}
