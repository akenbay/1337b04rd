package services

import (
	"context"
	"errors"
	"mime/multipart"
	"reflect"
	"testing"

	"1337b04rd/internal/domain"
)

// --------------------
// Mocks for PostService dependencies
// --------------------

type MockPostRepo struct {
	savedPost  *domain.Post
	savePost   *domain.Post
	saveErr    error
	findPost   *domain.Post
	findErr    error
	active     []*domain.Post
	archived   []*domain.Post
	archiveErr error
}

func (m *MockPostRepo) Save(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	m.savedPost = post
	return m.savePost, m.saveErr
}

func (m *MockPostRepo) FindByID(ctx context.Context, id string) (*domain.Post, error) {
	return m.findPost, m.findErr
}

func (m *MockPostRepo) FindActive(ctx context.Context) ([]*domain.Post, error) {
	return m.active, m.findErr
}

func (m *MockPostRepo) FindArchived(ctx context.Context) ([]*domain.Post, error) {
	return m.archived, m.findErr
}

func (m *MockPostRepo) ArchiveOldPosts(ctx context.Context) error {
	return m.archiveErr
}

// --------------------
// Tests
// --------------------

func TestCreatePost_Success(t *testing.T) {
	mockRepo := &MockPostRepo{savePost: &domain.Post{Title: "ok"}}
	mockImageStorage := &MockImageStorage{storeURL: "http://img.url"}
	mockFileUtils := &MockFileUtils{bytes: []byte("img")}
	mockUserRepo := &MockUserRepo{findUser: &domain.User{SessionID: "u1", Username: "Test User"}}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewPostService(mockRepo, mockImageStorage, mockFileUtils, realUserService, "bucket123")

	req := &domain.CreatePostReq{
		Title:     "Post title",
		Content:   "Post content",
		SessionID: "u1",
		ImageData: []*multipart.FileHeader{{Filename: "img1.png"}},
	}

	post, err := svc.CreatePost(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Title != "ok" {
		t.Errorf("expected returned post Title 'ok', got '%s'", post.Title)
	}
	if mockRepo.savedPost.User.SessionID != "u1" {
		t.Errorf("expected user with session 'u1', got %#v", mockRepo.savedPost.User)
	}
	if len(mockRepo.savedPost.ImageURLs) != 1 || mockRepo.savedPost.ImageURLs[0] != "http://img.url" {
		t.Errorf("expected image URL 'http://img.url', got %#v", mockRepo.savedPost.ImageURLs)
	}
}

func TestCreatePost_ValidateImageFails(t *testing.T) {
	mockRepo := &MockPostRepo{}
	mockImageStorage := &MockImageStorage{}
	mockFileUtils := &MockFileUtils{validateErr: errors.New("invalid image")}
	mockUserRepo := &MockUserRepo{}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewPostService(mockRepo, mockImageStorage, mockFileUtils, realUserService, "bucket123")

	req := &domain.CreatePostReq{ImageData: []*multipart.FileHeader{{Filename: "img.png"}}}

	_, err := svc.CreatePost(context.Background(), req)
	if err == nil || err.Error() != "invalid image" {
		t.Fatalf("expected 'invalid image', got %v", err)
	}
}

func TestCreatePost_FileHeaderToBytesFails(t *testing.T) {
	mockRepo := &MockPostRepo{}
	mockImageStorage := &MockImageStorage{}
	mockFileUtils := &MockFileUtils{bytesErr: errors.New("bad bytes")}
	mockUserRepo := &MockUserRepo{}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewPostService(mockRepo, mockImageStorage, mockFileUtils, realUserService, "bucket123")

	req := &domain.CreatePostReq{ImageData: []*multipart.FileHeader{{Filename: "img.png"}}}

	_, err := svc.CreatePost(context.Background(), req)
	if err == nil || err.Error() != "bad bytes" {
		t.Fatalf("expected 'bad bytes', got %v", err)
	}
}

func TestCreatePost_StoreFails(t *testing.T) {
	mockRepo := &MockPostRepo{}
	mockImageStorage := &MockImageStorage{storeErr: errors.New("store fail")}
	mockFileUtils := &MockFileUtils{bytes: []byte("ok")}
	mockUserRepo := &MockUserRepo{}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewPostService(mockRepo, mockImageStorage, mockFileUtils, realUserService, "bucket123")

	req := &domain.CreatePostReq{ImageData: []*multipart.FileHeader{{Filename: "img.png"}}}

	_, err := svc.CreatePost(context.Background(), req)
	if err == nil || err.Error() != "store fail" {
		t.Fatalf("expected 'store fail', got %v", err)
	}
}

func TestCreatePost_FindUserFails(t *testing.T) {
	mockRepo := &MockPostRepo{}
	mockImageStorage := &MockImageStorage{storeURL: "http://img.url"}
	mockFileUtils := &MockFileUtils{bytes: []byte("ok")}
	mockUserRepo := &MockUserRepo{findErr: errors.New("no user")}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewPostService(mockRepo, mockImageStorage, mockFileUtils, realUserService, "bucket123")

	req := &domain.CreatePostReq{SessionID: "u1", ImageData: []*multipart.FileHeader{{Filename: "img.png"}}}

	_, err := svc.CreatePost(context.Background(), req)
	if err == nil || err.Error() != "no user" {
		t.Fatalf("expected 'no user', got %v", err)
	}
}

func TestGetPostByID_Success(t *testing.T) {
	expected := &domain.Post{Title: "test"}
	mockRepo := &MockPostRepo{findPost: expected}

	svc := NewPostService(mockRepo, nil, nil, UserService{}, "")

	got, err := svc.GetPostByID(context.Background(), "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}

func TestGetActivePosts_Success(t *testing.T) {
	expected := []*domain.Post{{Title: "p1"}}
	mockRepo := &MockPostRepo{active: expected}

	svc := NewPostService(mockRepo, nil, nil, UserService{}, "")

	got, err := svc.GetActivePosts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}

func TestGetArchivedPosts_Success(t *testing.T) {
	expected := []*domain.Post{{Title: "archived"}}
	mockRepo := &MockPostRepo{archived: expected}

	svc := NewPostService(mockRepo, nil, nil, UserService{}, "")

	got, err := svc.GetArchivedPosts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}

func TestArchivePosts_Success(t *testing.T) {
	mockRepo := &MockPostRepo{}

	svc := NewPostService(mockRepo, nil, nil, UserService{}, "")

	if err := svc.ArchivePosts(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestArchivePosts_Error(t *testing.T) {
	mockRepo := &MockPostRepo{archiveErr: errors.New("archive fail")}

	svc := NewPostService(mockRepo, nil, nil, UserService{}, "")

	if err := svc.ArchivePosts(context.Background()); err == nil || err.Error() != "archive fail" {
		t.Fatalf("expected 'archive fail', got %v", err)
	}
}
