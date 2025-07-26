package postgres

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
)

type CommentRepository struct {
	db *sql.DB
}

var _ domain.CommentRepository = (*CommentRepository)(nil)

func NewCommentRepository(db *sql.DB, defaultBucket string) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Save(ctx context.Context, comment *domain.Comment) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO comments (
            post_id, parent_id, content, 
            image_key, bucket_name
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING comment_id, created_at
    `

	var imageKey, bucketName sql.NullString
	if comment.ImageKey != nil {
		imageKey = sql.NullString{String: *comment.ImageKey, Valid: true}
		bucketName = sql.NullString{String: *comment.BucketName, Valid: true}
	}

	err = tx.QueryRowContext(ctx, query,
		comment.PostID,
		comment.ParentID,
		comment.Content,
		imageKey,
		bucketName,
	).Scan(
		&comment.ID,        // Populate the generated UUID
		&comment.CreatedAt, // Get actual DB timestamp
	)

	if err != nil {
		return "", err
	}

	return comment.ID, tx.Commit()
}

func (r *CommentRepository) FindByPostID(ctx context.Context, postid string) ([]*domain.Comment, error) {
	query := `
		SELECT 
			c.comment_id, c.post_id, c.parent_id,
			c.content, c.image_key, c.bucket_name,
			c.created_at,
			u.session_id, u.avatar_url,
			u.username
		FROM comments c
		JOIN user_sessions u ON c.session_id = u.session_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, postid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var comment domain.Comment
		var imageKey, bucketName sql.NullString

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Content,
			&imageKey,
			&bucketName,
			&comment.CreatedAt,
			&comment.User.SessionID,
			&comment.User.AvatarURL,
			&comment.User.Username,
		)
		if err != nil {
			return nil, err
		}

		if imageKey.Valid {
			comment.ImageKey = &imageKey.String
			comment.BucketName = &bucketName.String
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}
