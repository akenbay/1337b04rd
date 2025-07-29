package services

import (
	"1337b04rd/internal/domain"
	"context"
	"log/slog"
)

type CommentService struct {
	commentRepo   domain.CommentRepository
	userService   UserService
	imageStorage  domain.ImageStorageAPI
	fileUtils     domain.FileUtils
	defaultBucket string
}

func NewCommentService(commentRepo domain.CommentRepository, userService UserService, imageStorage domain.ImageStorageAPI, fileUtils domain.FileUtils, defaultBucket string) *CommentService {
	return &CommentService{
		commentRepo:   commentRepo,
		userService:   userService,
		imageStorage:  imageStorage,
		fileUtils:     fileUtils,
		defaultBucket: defaultBucket,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, createCommentReq *domain.CreateCommentReq) (string, error) {
	var comment domain.Comment

	for _, fileheader := range createCommentReq.ImageData {

		// Validate if its image
		if err := s.fileUtils.ValidateImage(fileheader); err != nil {
			slog.Error("Failed to validate the image", "error", err)
			return "", err
		}

		// Convert into bytes
		fileBytes, err := s.fileUtils.FileHeaderToBytes(fileheader)
		if err != nil {
			slog.Error("Failed to convert image into bytes.")
			return "", err
		}

		imageURL, err := s.imageStorage.Store(fileBytes, s.defaultBucket)
		if err != nil {
			slog.Error("Failed to store the image", "error", err)
			return "", err
		}
		comment.ImageURLs = append(comment.ImageURLs, imageURL)
	}

	comment.Content = createCommentReq.Content
	sessionID := createCommentReq.SessionID

	user, err := s.userService.FindUserByID(ctx, sessionID)
	if err != nil {
		return "", err
	}

	comment.User = *user

	return s.commentRepo.Save(ctx, &comment)
}

func (s *CommentService) LoadComments(ctx context.Context, postid string) ([]*domain.Comment, error) {
	return s.commentRepo.FindByPostID(ctx, postid)
}
