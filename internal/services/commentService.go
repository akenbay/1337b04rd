package services

import (
	"1337b04rd/internal/domain"
	"1337b04rd/pkg/logger"
	"context"
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
	logger.Info("Service create comment:")

	var comment domain.Comment

	for _, fileheader := range createCommentReq.ImageData {

		// Validate if its image
		if err := s.fileUtils.ValidateImage(fileheader); err != nil {
			logger.Error("Failed to validate the image", "error", err)
			return "", err
		}

		// Convert into bytes
		fileBytes, err := s.fileUtils.FileHeaderToBytes(fileheader)
		if err != nil {
			logger.Error("Failed to convert image into bytes.")
			return "", err
		}

		imageURL, err := s.imageStorage.Store(fileBytes, s.defaultBucket)
		if err != nil {
			logger.Error("Failed to store the image", "error", err)
			return "", err
		}
		comment.ImageURLs = append(comment.ImageURLs, imageURL)
	}

	logger.Info("Preccessed and stored images from comment")

	comment.Content = createCommentReq.Content
	comment.PostID = createCommentReq.PostID
	sessionID := createCommentReq.SessionID

	if createCommentReq.ParentID != nil {
		comment.ParentID = createCommentReq.ParentID
	}

	user, err := s.userService.FindUserByID(ctx, sessionID)
	if err != nil {
		return "", err
	}

	comment.User = *user

	logger.Info("Found user by ID and assigned it to comment")

	return s.commentRepo.Save(ctx, &comment)
}

func (s *CommentService) LoadComments(ctx context.Context, postid string) ([]*domain.Comment, error) {
	return s.commentRepo.FindByPostID(ctx, postid)
}
