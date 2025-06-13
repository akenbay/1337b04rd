package services

import (
	"1337b04rd/internal/domain"
	"context"
)

type PostService struct {
	postRepo      domain.PostRepository
	userService   UserService
	imageStorage  domain.ImageStorageAPI
	defaultBucket string
}

func NewPostService(postRepo domain.PostRepository, imageStorage domain.ImageStorageAPI) *PostService {
	return &PostService{
		postRepo:     postRepo,
		imageStorage: imageStorage,
	}
}

func (s *PostService) CreatePost(ctx context.Context, createPostReq *domain.CreatePostReq) error {
	var post domain.Post

	if createPostReq.ImageData == nil {
	}

	image_key, err := s.imageStorage.Store(createPostReq.ImageData, s.defaultBucket)
	if err != nil {
		return err
	}

	post.Title = createPostReq.Title
	post.Content = createPostReq.Content
	post.ImageKey = &image_key
	post.BucketName = &s.defaultBucket
	sessionID := createPostReq.SessionID

	user, err := s.userService.FindUserByID(ctx, sessionID)
	if err != nil {
		return err
	}

	post.User = *user

	if err := post.Validate(); err != nil {
		return err
	}

	return s.postRepo.Save(ctx, &post)
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
