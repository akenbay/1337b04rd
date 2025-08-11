package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"1337b04rd/internal/domain"

	"github.com/lib/pq"
)

type PostRepository struct {
	db *sql.DB
}

var _ domain.PostRepository = (*PostRepository)(nil)

func NewPostRepository(db *sql.DB, defaultBucket string) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	if err := post.Validate(); err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO posts (
            session_id, title, content, 
            image_urls
        ) VALUES ($1, $2, $3, $4)
        RETURNING post_id, created_at, updated_at
    `

	err = tx.QueryRowContext(ctx, query,
		post.User.SessionID,
		post.Title,
		post.Content,
		pq.Array(post.ImageURLs),
	).Scan(
		&post.ID,        // Populate the generated UUID
		&post.CreatedAt, // Get actual DB timestamp
		&post.UpdatedAt, // Get actual DB timestamp
	)
	if err != nil {
		return nil, err
	}

	return post, tx.Commit()
}

func (r *PostRepository) FindByID(ctx context.Context, id string) (*domain.Post, error) {
	query := `
		SELECT 
			p.post_id, p.title, p.content,
			p.image_urls,
			p.created_at, p.updated_at, p.is_archived, p.archived_at,
			u.session_id, u.avatar_url, 
			u.username
		FROM posts p
		JOIN user_sessions u ON p.session_id = u.session_id
		WHERE p.post_id = $1
	`

	var post domain.Post
	var imageURLs pq.StringArray
	var archivedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&imageURLs,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsArchived,
		&archivedAt,
		&post.User.SessionID,
		&post.User.AvatarURL,
		&post.User.Username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	// Handle nullable fields
	if archivedAt.Valid {
		post.ArchivedAt = &archivedAt.Time
	}

	// Convert pq string array into massive
	post.ImageURLs = []string(imageURLs)

	return &post, nil
}

func (r *PostRepository) FindActive(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT 
			p.post_id, p.title, p.content,
			p.image_urls,
			p.created_at, p.updated_at,
			u.session_id, u.avatar_url,
			u.username
		FROM posts p
		JOIN user_sessions u ON p.session_id = u.session_id
		WHERE p.is_archived = FALSE
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		var imageURLs pq.StringArray

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&imageURLs,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.SessionID,
			&post.User.AvatarURL,
			&post.User.Username,
		)
		if err != nil {
			return nil, err
		}

		post.ImageURLs = []string(imageURLs)

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) FindArchived(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT 
			p.post_id, p.title, p.content,
			p.image_urls,
			p.created_at, p.updated_at,
			u.session_id, u.avatar_url,
			u.username
		FROM posts p
		JOIN user_sessions u ON p.session_id = u.session_id
		WHERE p.is_archived = TRUE
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var post domain.Post
		var imageURLs pq.StringArray

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&imageURLs,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.SessionID,
			&post.User.AvatarURL,
			&post.User.Username,
		)
		if err != nil {
			return nil, err
		}

		post.ImageURLs = []string(imageURLs)

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) ArchiveOldPosts(ctx context.Context) error {
	// Use the database function we defined in init.sql
	_, err := r.db.ExecContext(ctx, "SELECT archive_old_posts()")
	return err
}

func sqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
