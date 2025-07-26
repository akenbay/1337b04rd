package domain

import (
	"context"
	"time"
)

type Comment struct {
	ID         string  // UUID
	PostID     string  // Parent post UUID
	ParentID   *string // Nullable (for nested comments)
	User       User    // Embedded or reference SessionID
	Content    string
	ImageKey   *string // S3 object key (nullable)
	BucketName *string // S3 bucket (nullable)
	CreatedAt  time.Time
}

type CreateCommentReq struct {
	SessionID string
	Content   string
	ImageData []byte
}

type CommentRepository interface {
	Save(ctx context.Context, comment *Comment) (string, error)
	FindByPostID(ctx context.Context, postid string) ([]*Comment, error)
}
