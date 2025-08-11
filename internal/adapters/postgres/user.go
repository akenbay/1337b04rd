package postgres

import (
	"1337b04rd/internal/domain"
	"1337b04rd/pkg/logger"
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
	logger.Info("Save func")
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

	logger.Info("Created query and transaction")

	err = tx.QueryRowContext(ctx, query,
		avatarURL,
		name,
	).Scan(
		&user.SessionID, // Populate the generated UUID
	)

	logger.Info("Pushed the query into the sql")

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
		logger.Error("Postgres, error when changing username:", "error", err)
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetNumberOfUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM user_sessions").Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *UserRepository) FindByID(ctx context.Context, session_id string) (*domain.User, error) {
	query := `SELECT 
	session_id,
	avatar_url,
    username,
    created_at,
    expires_at 
	FROM user_sessions
	WHERE session_id=$1
	`

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, session_id).Scan(
		&user.SessionID,
		&user.AvatarURL,
		&user.Username,
		&user.CreatedAt,
		&user.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
