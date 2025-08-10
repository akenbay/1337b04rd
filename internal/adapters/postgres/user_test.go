package postgres_test

import (
	"1337b04rd/internal/adapters/postgres"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	// Adjust connection string according to your docker-compose service
	testDB, err = sql.Open("postgres", "postgres://latte:latte@db:5432/1337b04rd?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = testDB.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	// Make sure schema exists
	_, err = testDB.Exec(`
        CREATE TABLE IF NOT EXISTS user_sessions (
            session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            avatar_url TEXT,
            username TEXT,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            expires_at TIMESTAMPTZ DEFAULT NOW() + INTERVAL '1 day'
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestSaveAndFindByID(t *testing.T) {
	repo := postgres.NewUserRepository(testDB)
	ctx := context.Background()

	sessionID, err := repo.Save(ctx, "http://avatar.com/me.png", "test_user")
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if sessionID == "" {
		t.Fatal("expected non-empty sessionID")
	}

	user, err := repo.FindByID(ctx, sessionID)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if user.Username != "test_user" {
		t.Errorf("expected username 'test_user', got '%s'", user.Username)
	}
}

func TestChangeName(t *testing.T) {
	repo := postgres.NewUserRepository(testDB)
	ctx := context.Background()

	sessionID, _ := repo.Save(ctx, "http://avatar.com/2.png", "old_name")
	err := repo.ChangeName(ctx, "new_name", sessionID)
	if err != nil {
		t.Fatalf("ChangeName() error: %v", err)
	}

	user, _ := repo.FindByID(ctx, sessionID)
	if user.Username != "new_name" {
		t.Errorf("expected username 'new_name', got '%s'", user.Username)
	}
}

func TestGetNumberOfUsers(t *testing.T) {
	repo := postgres.NewUserRepository(testDB)
	ctx := context.Background()

	count, err := repo.GetNumberOfUsers(ctx)
	if err != nil {
		t.Fatalf("GetNumberOfUsers() error: %v", err)
	}
	if count <= 0 {
		t.Errorf("expected count > 0, got %d", count)
	}
}
