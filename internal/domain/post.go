package domain

import (
	"context"
	"errors"
	"time"
)

type Post struct {
	ID         string // UUID
	User       User   // Embedded or reference SessionID
	Title      string
	Content    string
	ImageKey   *string // S3 object key (nullable)
	BucketName *string // S3 bucket (nullable)
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsArchived bool
	ArchivedAt *time.Time
}

type CreatePostReq struct {
	SessionID string
	Title     string
	Content   string
	ImageData []byte
}

type PostRepository interface {
	Save(ctx context.Context, post *Post) (string, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	FindActive(ctx context.Context) ([]*Post, error)
	FindArchived(ctx context.Context) ([]*Post, error)
	ArchiveOldPosts(ctx context.Context) error
}

func (p *Post) Validate() error {
	if len(p.Title) < 5 {
		return errors.New("title too short")
	}
	return nil
}
