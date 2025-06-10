package services

import (
	"1337b04rd/internal/domain"
	"context"
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

func (s *PostService) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	return s.postRepo.FindByID(ctx, id)
}

func (s *PostService) GetActivePosts(ctx context.Context) ([]*domain.Post, error) {
	return s.postRepo.FindActive(ctx)
}

func (s *PostService) GetArchivedPosts(ctx context.Context) ([]*domain.Post, error) {
	return s.postRepo.FindArchived(ctx)
}

func (s *PostService) ArchivePosts(ctx context.Context) error {
	return s.postRepo.ArchiveOldPosts(ctx)
}
