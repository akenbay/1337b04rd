package services

import (
	"context"

	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/repositories"
)

type PostService struct {
	postRepo repositories.PostRepository
}

func NewPostService(postRepo repositories.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(ctx context.Context, post *models.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	return s.postRepo.Save(ctx, post)
}
