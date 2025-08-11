package services

import (
	"1337b04rd/internal/domain"
	"1337b04rd/pkg/logger"
	"context"
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
	logger.Info("got number of users:", "count", count)
	if err != nil {
		return "", err
	}

	userOutlook, err := s.userOutlookAPI.GenerateAvatarAndName(count + 1)
	logger.Info("Generated avatar and username", "name", userOutlook.Name)
	if err != nil {
		logger.Error("Failed to generate avatar and name", "error", err)
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
