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

func (s *UserService) Save(ctx context.Context) error {
	count, err := s.userRepo.GetNumberOfUsers(ctx)
	if err != nil {
		return err
	}

	userOutlook, err := s.userOutlookAPI.GenerateAvatarAndName(count)
	if err != nil {
		return err
	}

	return s.userRepo.Save(ctx, userOutlook.AvatarURL, userOutlook.Name)
}
