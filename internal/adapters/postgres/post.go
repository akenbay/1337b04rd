package postgres

import (
	"1337b04rd/internal/domain/models"
	"1337b04rd/internal/ports/repositories"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostRepository struct {
	db            *sql.DB
	defaultBucket string // Your main S3 bucket name
}

var _ repositories.PostRepository = (*PostRepository)(nil)

func NewPostRepository(db *sql.DB, defaultBucket string) *PostRepository {
	return &PostRepository{
		db:            db,
		defaultBucket: defaultBucket,
	}
}

func (r *PostRepository) Save(ctx context.Context, post *models.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}

	// Use transaction for atomic operations
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Handle image reference
	var imageKey, bucketName sql.NullString
	if post.ImageKey != nil {
		imageKey = sql.NullString{String: *post.ImageKey, Valid: true}
		bucketName = sql.NullString{
			String: r.getBucketName(post.BucketName),
			Valid:  true,
		}
	}

	query := `
		INSERT INTO posts (
			id, user_id, title, content, 
			image_key, bucket_name, 
			created_at, is_archived, archived_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			image_key = EXCLUDED.image_key,
			bucket_name = EXCLUDED.bucket_name,
			is_archived = EXCLUDED.is_archived,
			archived_at = EXCLUDED.archived_at
	`
	_, err = tx.ExecContext(ctx, query,
		post.ID,
		post.User.ID,
		post.Title,
		post.Content,
		imageKey,
		bucketName,
		post.CreatedAt,
		post.IsArchived,
		sqlNullTime(post.ArchivedAt),
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// Helper to safely handle nil bucket names
func (r *PostRepository) getBucketName(bucket *string) string {
	if bucket != nil {
		return *bucket
	}
	return r.defaultBucket
}

// Helper for nullable time fields
func sqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func (r *PostRepository) FindByID(ctx context.Context, id string) (*models.Post, error) {
	query := `
		SELECT 
			p.id, p.title, p.content, 
			p.image_key, p.bucket_name,
			p.created_at, p.is_archived, p.archived_at,
			u.id, u.session_id, u.display_name, u.avatar_url
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var post models.Post
	var imageKey, bucketName sql.NullString
	var archivedAt sql.NullTime

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&imageKey,
		&bucketName,
		&post.CreatedAt,
		&post.IsArchived,
		&archivedAt,
		&post.User.ID,
		&post.User.SessionID,
		&post.User.DisplayName,
		&post.User.AvatarURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrPostNotFound
		}
		return nil, err
	}

	// Handle nullable fields
	if imageKey.Valid {
		post.ImageKey = &imageKey.String
		post.BucketName = &bucketName.String
	}
	if archivedAt.Valid {
		post.ArchivedAt = &archivedAt.Time
	}

	return &post, nil
}

func (r *PostRepository) FindActive(ctx context.Context) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.image_url, p.created_at, p.updated_at,
		       u.id, u.display_name, u.avatar_url
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.archived = false
		ORDER BY p.updated_at DESC
	`
	// Implementation similar to FindByID but with rows.Next()
}

func (r *PostRepository) ArchiveOldPosts(ctx context.Context) error {
	query := `
		UPDATE posts
		SET archived = true
		WHERE updated_at < NOW() - INTERVAL '15 minutes'
		AND id NOT IN (
			SELECT DISTINCT post_id FROM comments
			WHERE created_at > NOW() - INTERVAL '15 minutes'
		)
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
