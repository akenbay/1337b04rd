package postgres

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostRepository struct {
	db            *sql.DB
	defaultBucket string
}

var _ domain.PostRepository = (*PostRepository)(nil)

func NewPostRepository(db *sql.DB, defaultBucket string) *PostRepository {
	return &PostRepository{
		db:            db,
		defaultBucket: defaultBucket,
	}
}

func (r *PostRepository) Save(ctx context.Context, post *domain.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO posts (
            session_id, title, content, 
            image_key, bucket_name
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING post_id, created_at, updated_at
    `

	var imageKey, bucketName sql.NullString
	if post.ImageKey != nil {
		imageKey = sql.NullString{String: *post.ImageKey, Valid: true}
		bucketName = sql.NullString{String: r.getBucketName(post.BucketName), Valid: true}
	}

	err = tx.QueryRowContext(ctx, query,
		post.User.SessionID,
		post.Title,
		post.Content,
		imageKey,
		bucketName,
	).Scan(
		&post.ID,        // Populate the generated UUID
		&post.CreatedAt, // Get actual DB timestamp
		&post.UpdatedAt, // Get actual DB timestamp
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostRepository) FindByID(ctx context.Context, id string) (*domain.Post, error) {
	query := `
		SELECT 
			p.post_id, p.title, p.content,
			p.image_key, p.bucket_name,
			p.created_at, p.updated_at, p.is_archived, p.archived_at,
			u.session_id, u.avatar_url, 
			COALESCE(u.custom_name, u.character_name) as display_name
		FROM posts p
		JOIN user_sessions u ON p.session_id = u.session_id
		WHERE p.post_id = $1 AND p.is_archived = FALSE
	`

	var post domain.Post
	var imageKey, bucketName sql.NullString
	var archivedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&imageKey,
		&bucketName,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsArchived,
		&archivedAt,
		&post.User.SessionID,
		&post.User.AvatarURL,
		&post.User.CharacterName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
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

func (r *PostRepository) FindActive(ctx context.Context) ([]*domain.Post, error) {
	query := `
		SELECT 
			p.post_id, p.title, p.content,
			p.image_key, p.bucket_name,
			p.created_at, p.updated_at,
			u.session_id, u.avatar_url,
			COALESCE(u.custom_name, u.character_name) as display_name
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
		var imageKey, bucketName sql.NullString

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&imageKey,
			&bucketName,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.SessionID,
			&post.User.AvatarURL,
			&post.User.CharacterName,
		)
		if err != nil {
			return nil, err
		}

		if imageKey.Valid {
			post.ImageKey = &imageKey.String
			post.BucketName = &bucketName.String
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *PostRepository) ArchiveOldPosts(ctx context.Context) error {
	// Use the database function we defined in init.sql
	_, err := r.db.ExecContext(ctx, "SELECT archive_old_posts()")
	return err
}

// Helper functions
func (r *PostRepository) getBucketName(bucket *string) string {
	if bucket != nil {
		return *bucket
	}
	return r.defaultBucket
}

func sqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
