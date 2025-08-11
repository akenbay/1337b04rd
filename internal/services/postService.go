package services

import (
	"1337b04rd/internal/domain"
	"1337b04rd/pkg/logger"
	"context"
)

type PostService struct {
	postRepo      domain.PostRepository
	userService   UserService
	imageStorage  domain.ImageStorageAPI
	fileUtils     domain.FileUtils
	defaultBucket string
}

func NewPostService(postRepo domain.PostRepository, imageStorage domain.ImageStorageAPI, fileUtils domain.FileUtils, userService UserService, defaultBucket string) *PostService {
	return &PostService{
		postRepo:      postRepo,
		imageStorage:  imageStorage,
		fileUtils:     fileUtils,
		userService:   userService,
		defaultBucket: defaultBucket,
	}
}

func (s *PostService) CreatePost(ctx context.Context, createPostReq *domain.CreatePostReq) (*domain.Post, error) {
	var post domain.Post

	for _, fileheader := range createPostReq.ImageData {

		// Validate if its image
		if err := s.fileUtils.ValidateImage(fileheader); err != nil {
			logger.Error("Failed to validate the image", "error", err)
			return nil, err
		}

		// Convert into bytes
		fileBytes, err := s.fileUtils.FileHeaderToBytes(fileheader)
		if err != nil {
			logger.Error("Failed to convert image into bytes.")
			return nil, err
		}

		imageURL, err := s.imageStorage.Store(fileBytes, s.defaultBucket)
		if err != nil {
			logger.Error("Failed to store the image:", "error", err)
			return nil, err
		}
		post.ImageURLs = append(post.ImageURLs, imageURL)
	}

	post.Title = createPostReq.Title
	post.Content = createPostReq.Content
	sessionID := createPostReq.SessionID

	user, err := s.userService.FindUserByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	post.User = *user

	if err := post.Validate(); err != nil {
		return nil, err
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
