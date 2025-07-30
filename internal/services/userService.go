package services

import (
	"1337b04rd/internal/domain"
	"context"
	"log/slog"
)

type UserService struct {
	userRepo       domain.UserRepository
	userOutlookAPI domain.UserOutlookAPI
}

func NewUserService(userRepo domain.UserRepository, userOutlookAPI domain.UserOutlookAPI) *UserService {
	return &UserService{
		userRepo:       userRepo,
		userOutlookAPI: userOutlookAPI,
	}
}

func (s *UserService) CreateUserAndGetID(ctx context.Context) (string, error) {
	count, err := s.userRepo.GetNumberOfUsers(ctx)
	slog.Info("got number of users")
	if err != nil {
		return "", err
	}

	userOutlook, err := s.userOutlookAPI.GenerateAvatarAndName(count)
	slog.Info("Generated avatar and username")
	if err != nil {
		slog.Error("Failed to generate avatar and name", "error", err)
		return "", err
	}

	return s.userRepo.Save(ctx, userOutlook.AvatarURL, userOutlook.Name)
}

func (s *UserService) ChangeUsername(ctx context.Context, session_id string, newUsername string) error {
	return s.userRepo.ChangeName(ctx, newUsername, session_id)
}

func (s *UserService) FindUserByID(ctx context.Context, session_id string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, session_id)
}
