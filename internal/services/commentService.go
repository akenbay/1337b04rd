package services

import (
	"1337b04rd/internal/domain"
	"context"
	"log/slog"
)

type CommentService struct {
	commentRepo    domain.CommentRepository
	userService    UserService
	imageStorage   domain.ImageStorageAPI
	imageValidator domain.ImageValidator
	defaultBucket  string
}

func (s *CommentService) CreateComment(ctx context.Context, createCommentReq *domain.CreateCommentReq) (string, error) {
	var comment domain.Comment

	if err := s.imageValidator.Validate(createCommentReq.ImageData); err != nil {
		slog.Error("Failed to validate the image", "error", err)
		return "", err
	}

	if createCommentReq.ImageData != nil {
		image_key, err := s.imageStorage.Store(createCommentReq.ImageData, s.defaultBucket)
		if err != nil {
			slog.Error("Failed to store the image", "error", err)
			return "", err
		}
		comment.ImageKey = &image_key
	} else {
		comment.ImageKey = nil
	}

	comment.Content = createCommentReq.Content
	comment.BucketName = &s.defaultBucket
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
