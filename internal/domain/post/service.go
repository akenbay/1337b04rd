package post

import (
	"context"
	"time"
)

type Service struct {
	repo         Repository
	imageStorage ImageStorage
	avatarSvc    AvatarService
}

func NewService(repo Repository, imageStorage ImageStorage, avatarSvc AvatarService) *Service {
	return &Service{
		repo:         repo,
		imageStorage: imageStorage,
		avatarSvc:    avatarSvc,
	}
}

func (s *Service) CreatePost(ctx context.Context, title, content string, imageData []byte, sessionID string) (*Post, error) {
	if title == "" {
		return nil, ErrPostTitleRequired
	}
	if content == "" {
		return nil, ErrPostContentRequired
	}

	user, err := s.avatarSvc.GetUser(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	post := &Post{
		ID:        generateID(),
		Title:     title,
		Content:   content,
		User:      *user,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if len(imageData) > 0 {
		imageURL, err := s.imageStorage.StoreImage(ctx, imageData)
		if err != nil {
			return nil, err
		}
		post.ImageURL = imageURL
	}

	if err := s.repo.Save(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}
