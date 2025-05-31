package repositories

import (
	"1337b04rd/internal/domain/post"
	"context"
)

// Port definition - how our app wants to talk to storage
type PostRepository interface {
	Save(ctx context.Context, p *post.Post) error
	FindByID(ctx context.Context, id string) (*post.Post, error)
}
