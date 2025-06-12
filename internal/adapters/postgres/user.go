package postgres

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

var _ domain.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, avatarURL string, name string) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO user_sessions (
            avatar_url, username
        ) VALUES ($1, $2)
        RETURNING session_id
    `

	var user domain.User

	err = tx.QueryRowContext(ctx, query,
		avatarURL,
		name,
	).Scan(
		&user.SessionID, // Populate the generated UUID
	)

	if err != nil {
		return "", err
	}

	return user.SessionID, tx.Commit()
}

func (r *UserRepository) ChangeName(ctx context.Context, newName string, sessionID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        UPDATE user_sessions
		SET username = $1
		WHERE session_id = $2
		RETURNING 
			session_id,
            avatar_url,
            username,
            created_at,
            expires_at
    `

	var user domain.User

	err = tx.QueryRowContext(ctx, query,
		newName,
		sessionID,
	).Scan(
		&user.SessionID,
		&user.AvatarURL,
		&user.Username,
		&user.CreatedAt,
		&user.ExpiresAt,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetNumberOfUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM your_table").Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}
