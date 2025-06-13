package services

import (
	"1337b04rd/internal/domain"
	"context"
)

type UserService struct {
	userRepo       domain.UserRepository
	userOutlookAPI domain.UserOutlookAPI
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUserAndGetID(ctx context.Context) (string, error) {
	count, err := s.userRepo.GetNumberOfUsers(ctx)
	if err != nil {
		return "", err
	}

	userOutlook, err := s.userOutlookAPI.GenerateAvatarAndName(count)
	if err != nil {
		return "", err
	}

	return s.userRepo.Save(ctx, userOutlook.AvatarURL, userOutlook.Name)
}

func (s *UserService) FindUserByID(ctx context.Context, session_id string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, session_id)
}
