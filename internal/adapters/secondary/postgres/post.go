package postgres

import (
	"1337b04rd/ports/repositories"
	"context"
	"database/sql"
)

// Adapter implementation
type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) repositories.PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Save(ctx context.Context, p *post.Post) error {
	// Actual PostgreSQL implementation
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO posts (...) VALUES (...)",
		p.ID, p.Title /* etc */)
	return err
}
