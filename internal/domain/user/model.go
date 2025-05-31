package user

import "time"

// User represents an anonymous board user
type User struct {
	SessionID   string    `json:"-"`          // Never exposed to clients
	DisplayName string    `json:"name"`       // Can be changed by user
	AvatarURL   string    `json:"avatar_url"` // From Rick&Morty API
	CreatedAt   time.Time `json:"-"`          // Internal use only
}

// NewUser creates a user with required fields
func NewUser(sessionID, name, avatarURL string) *User {
	return &User{
		SessionID:   sessionID,
		DisplayName: name,
		AvatarURL:   avatarURL,
		CreatedAt:   time.Now(),
	}
}

// IsValid checks required fields
func (u *User) IsValid() bool {
	return u.SessionID != "" && u.AvatarURL != "" && u.DisplayName != ""
}
