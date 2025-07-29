package domain

import (
	"context"
	"errors"
	"time"
)

type Post struct {
	ID         string // UUID which generates in SQL itself
	User       User   // Embedded or reference SessionID
	Title      string
	Content    string
	ImageURLs  []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsArchived bool
	ArchivedAt *time.Time
}

// Structure for creating post request

type CreatePostReq struct {
	SessionID  string
	Title      string
	Content    string
	ImageDatas [][]byte
}

// Description of the functions that manipulate the database

type PostRepository interface {
	Save(ctx context.Context, post *Post) (string, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	FindActive(ctx context.Context) ([]*Post, error)
	FindArchived(ctx context.Context) ([]*Post, error)
	ArchiveOldPosts(ctx context.Context) error
}

// Validation of title length

func (p *Post) Validate() error {
	if len(p.Title) < 5 {
		return errors.New("title too short")
	}
	return nil
}
