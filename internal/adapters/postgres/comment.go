package postgres

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"

	"github.com/lib/pq"
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
            image_urls
        ) VALUES ($1, $2, $3, $4)
        RETURNING comment_id, created_at
    `

	err = tx.QueryRowContext(ctx, query,
		comment.PostID,
		comment.ParentID,
		comment.Content,
		pq.Array(comment.ImageURLs),
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
			c.content, c.image_urls,
			c.created_at,
			u.session_id, u.avatar_url,
			u.username
		FROM comments c
		LEFT JOIN user_sessions u ON c.session_id = u.session_id
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
		var imageURLs pq.StringArray

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Content,
			&imageURLs,
			&comment.CreatedAt,
			&comment.User.SessionID,
			&comment.User.AvatarURL,
			&comment.User.Username,
		)
		if err != nil {
			return nil, err
		}

		comment.ImageURLs = []string(imageURLs)

		comments = append(comments, &comment)
	}

	return comments, nil
}
