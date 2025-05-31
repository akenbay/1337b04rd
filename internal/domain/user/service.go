package user

import (
	"context"
)

// Service handles user-related operations
type Service struct {
	repo      Repository
	avatarAPI AvatarService
}

// AvatarService defines how we get avatars (secondary port)
type AvatarService interface {
	GetRandomAvatar(ctx context.Context) (name string, url string, err error)
}

func NewService(repo Repository, avatarAPI AvatarService) *Service {
	return &Service{
		repo:      repo,
		avatarAPI: avatarAPI,
	}
}

// GetOrCreateUser handles user session management
func (s *Service) GetOrCreateUser(ctx context.Context, sessionID string) (*User, error) {
	// Try to get existing user
	existing, err := s.repo.FindBySessionID(ctx, sessionID)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new user with avatar
	name, avatarURL, err := s.avatarAPI.GetRandomAvatar(ctx)
	if err != nil {
		return nil, err
	}

	newUser := NewUser(sessionID, name, avatarURL)
	if err := s.repo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// UpdateDisplayName changes how user appears on posts
func (s *Service) UpdateDisplayName(ctx context.Context, sessionID, newName string) (*User, error) {
	if newName == "" {
		return nil, ErrInvalidName
	}

	user, err := s.repo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	user.DisplayName = newName
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
