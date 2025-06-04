package models

import (
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

func (p *Post) Validate() error {
	if len(p.Title) < 5 {
		return errors.New("title too short")
	}
	return nil
}
