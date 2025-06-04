package repositories

import (
	"1337b04rd/internal/domain/models"
	"context"
)

type UserRepository interface {
	Save(ctx context.Context, post *models.Post) error
	FindByID(ctx context.Context, id string) (*models.Post, error)
	FindActive(ctx context.Context) ([]*models.Post, error)
	ArchiveOldPosts(ctx context.Context) error
}
