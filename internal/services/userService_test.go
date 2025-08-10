package services_test

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"context"
	"errors"
	"testing"
	"time"
)

// --- Fake implementations ---

type fakeUserRepo struct {
	count       int
	saveID      string
	users       map[string]*domain.User
	saveErr     error
	countErr    error
	changeErr   error
	findByIDErr error
}

func (f *fakeUserRepo) ChangeName(ctx context.Context, newName string, sessionID string) error {
	return f.changeErr
}

func (f *fakeUserRepo) Save(ctx context.Context, avatarURL string, name string) (string, error) {
	if f.saveErr != nil {
		return "", f.saveErr
	}
	f.saveID = "generated-id"
	return f.saveID, nil
}

func (f *fakeUserRepo) GetNumberOfUsers(ctx context.Context) (int, error) {
	if f.countErr != nil {
		return 0, f.countErr
	}
	return f.count, nil
}

func (f *fakeUserRepo) FindByID(ctx context.Context, sessionID string) (*domain.User, error) {
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return f.users[sessionID], nil
}

type fakeOutlookAPI struct {
	avatar string
	name   string
	err    error
}

func (f *fakeOutlookAPI) GenerateAvatarAndName(userNumber int) (*domain.UserOutlook, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &domain.UserOutlook{
		AvatarURL: f.avatar,
		Name:      f.name,
	}, nil
}

// --- Tests ---

func TestCreateUserAndGetID_Success(t *testing.T) {
	repo := &fakeUserRepo{count: 3}
	api := &fakeOutlookAPI{avatar: "avatar.png", name: "Morty"}
	svc := services.NewUserService(repo, api)

	id, err := svc.CreateUserAndGetID(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != "generated-id" {
		t.Errorf("expected generated-id, got %s", id)
	}
}

func TestCreateUserAndGetID_RepoCountError(t *testing.T) {
	repo := &fakeUserRepo{countErr: errors.New("count fail")}
	api := &fakeOutlookAPI{}
	svc := services.NewUserService(repo, api)

	_, err := svc.CreateUserAndGetID(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestChangeUsername(t *testing.T) {
	repo := &fakeUserRepo{}
	api := &fakeOutlookAPI{}
	svc := services.NewUserService(repo, api)

	err := svc.ChangeUsername(context.Background(), "sid", "newname")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFindUserByID(t *testing.T) {
	repo := &fakeUserRepo{
		users: map[string]*domain.User{
			"sid": {
				SessionID: "sid",
				Username:  "Rick",
				CreatedAt: time.Now(),
			},
		},
	}
	api := &fakeOutlookAPI{}
	svc := services.NewUserService(repo, api)

	user, err := svc.FindUserByID(context.Background(), "sid")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != "Rick" {
		t.Errorf("expected Rick, got %s", user.Username)
	}
}
