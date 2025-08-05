package domain

import (
	"context"
	"mime/multipart"
	"time"
)

type Comment struct {
	ID        string  // UUID
	PostID    string  // Parent post UUID
	ParentID  *string // Nullable (for nested comments)
	User      User    // Embedded or reference SessionID
	Content   string
	ImageURLs []string
	CreatedAt time.Time
}

type CreateCommentReq struct {
	SessionID string
	PostID    string
	Content   string
	ParentID  *string
	ImageData []*multipart.FileHeader
}

type CommentRepository interface {
	Save(ctx context.Context, comment *Comment) (string, error)
	FindByPostID(ctx context.Context, postid string) ([]*Comment, error)
}
