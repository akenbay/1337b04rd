package domain

import (
	"context"
	"time"
)

type User struct {
	SessionID string // UUID
	AvatarURL string // From Rick & Morty API
	Username  string // From API, yet user can override
	CreatedAt time.Time
	ExpiresAt time.Time // Session expiry (1 week)
}

type UserRepository interface {
	ChangeName(ctx context.Context, newName string, sessionID string) error
	Save(ctx context.Context, avatarURL string, name string) error
	GetNumberOfUsers(ctx context.Context) (int, error)
}
