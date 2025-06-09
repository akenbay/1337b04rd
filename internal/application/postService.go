package services

import (
	"context"

	"1337b04rd/internal/domain"
)

type PostService struct {
	postRepo domain.PostRepository
}

func NewPostService(postRepo domain.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}
	return s.postRepo.Save(ctx, post)
}
