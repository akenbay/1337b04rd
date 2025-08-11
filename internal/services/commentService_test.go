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
// Mocks for CommentService dependencies
// --------------------

type MockCommentRepo struct {
	savedComment *domain.Comment
	saveID       string
	saveErr      error
	comments     []*domain.Comment
	findErr      error
}

func (m *MockCommentRepo) Save(ctx context.Context, comment *domain.Comment) (string, error) {
	m.savedComment = comment
	return m.saveID, m.saveErr
}

func (m *MockCommentRepo) FindByPostID(ctx context.Context, postid string) ([]*domain.Comment, error) {
	return m.comments, m.findErr
}

type MockImageStorage struct {
	storeURL string
	storeErr error
}

func (m *MockImageStorage) Store(data []byte, bucket string) (string, error) {
	return m.storeURL, m.storeErr
}

type MockFileUtils struct {
	validateErr error
	bytes       []byte
	bytesErr    error
}

func (m *MockFileUtils) ValidateImage(fileHeader *multipart.FileHeader) error {
	return m.validateErr
}

func (m *MockFileUtils) FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	return m.bytes, m.bytesErr
}

// --------------------
// Mocks for UserService dependencies
// --------------------

type MockUserRepo struct {
	findUser *domain.User
	findErr  error
}

func (m *MockUserRepo) GetNumberOfUsers(ctx context.Context) (int, error) { return 0, nil }
func (m *MockUserRepo) Save(ctx context.Context, avatarURL, name string) (string, error) {
	return "", nil
}
func (m *MockUserRepo) ChangeName(ctx context.Context, newName, sessionID string) error { return nil }
func (m *MockUserRepo) FindByID(ctx context.Context, sessionID string) (*domain.User, error) {
	return m.findUser, m.findErr
}

type MockUserOutlookAPI struct{}

func (m *MockUserOutlookAPI) GenerateAvatarAndName(count int) (*domain.UserOutlook, error) {
	return &domain.UserOutlook{Name: "Test", AvatarURL: "avatar"}, nil
}

// --------------------
// Tests
// --------------------

func TestCreateComment_Success(t *testing.T) {
	mockRepo := &MockCommentRepo{saveID: "123"}
	mockImageStorage := &MockImageStorage{storeURL: "http://image.url"}
	mockFileUtils := &MockFileUtils{bytes: []byte("image data")}
	mockUserRepo := &MockUserRepo{
		findUser: &domain.User{SessionID: "u1", Username: "Test User"},
	}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewCommentService(mockRepo, realUserService, mockImageStorage, mockFileUtils, "bucket123")

	req := &domain.CreateCommentReq{
		Content:   "Hello World",
		PostID:    "p1",
		SessionID: "u1",
		ImageData: []*multipart.FileHeader{{Filename: "file1.png"}},
	}

	id, err := svc.CreateComment(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "123" {
		t.Errorf("expected id 123, got %s", id)
	}
	if mockRepo.savedComment.Content != "Hello World" {
		t.Errorf("expected content 'Hello World', got '%s'", mockRepo.savedComment.Content)
	}
	if len(mockRepo.savedComment.ImageURLs) != 1 || mockRepo.savedComment.ImageURLs[0] != "http://image.url" {
		t.Errorf("image URLs not saved correctly: %#v", mockRepo.savedComment.ImageURLs)
	}
	if mockRepo.savedComment.User.SessionID != "u1" {
		t.Errorf("user not assigned correctly: %#v", mockRepo.savedComment.User)
	}
}

func TestCreateComment_ValidateImageFails(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	mockImageStorage := &MockImageStorage{}
	mockFileUtils := &MockFileUtils{validateErr: errors.New("invalid image")}
	mockUserRepo := &MockUserRepo{}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewCommentService(mockRepo, realUserService, mockImageStorage, mockFileUtils, "bucket123")

	req := &domain.CreateCommentReq{
		ImageData: []*multipart.FileHeader{{Filename: "file1.png"}},
	}

	_, err := svc.CreateComment(context.Background(), req)
	if err == nil || err.Error() != "invalid image" {
		t.Fatalf("expected 'invalid image' error, got %v", err)
	}
}

func TestCreateComment_FileHeaderToBytesFails(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	mockImageStorage := &MockImageStorage{}
	mockFileUtils := &MockFileUtils{bytesErr: errors.New("cannot convert to bytes")}
	mockUserRepo := &MockUserRepo{}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewCommentService(mockRepo, realUserService, mockImageStorage, mockFileUtils, "bucket123")

	req := &domain.CreateCommentReq{
		ImageData: []*multipart.FileHeader{{Filename: "file1.png"}},
	}

	_, err := svc.CreateComment(context.Background(), req)
	if err == nil || err.Error() != "cannot convert to bytes" {
		t.Fatalf("expected 'cannot convert to bytes', got %v", err)
	}
}

func TestCreateComment_StoreFails(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	mockImageStorage := &MockImageStorage{storeErr: errors.New("store failed")}
	mockFileUtils := &MockFileUtils{bytes: []byte("image data")}
	mockUserRepo := &MockUserRepo{}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewCommentService(mockRepo, realUserService, mockImageStorage, mockFileUtils, "bucket123")

	req := &domain.CreateCommentReq{
		ImageData: []*multipart.FileHeader{{Filename: "file1.png"}},
	}

	_, err := svc.CreateComment(context.Background(), req)
	if err == nil || err.Error() != "store failed" {
		t.Fatalf("expected 'store failed', got %v", err)
	}
}

func TestCreateComment_FindUserFails(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	mockImageStorage := &MockImageStorage{storeURL: "http://image.url"}
	mockFileUtils := &MockFileUtils{bytes: []byte("image data")}
	mockUserRepo := &MockUserRepo{findErr: errors.New("user not found")}
	mockOutlook := &MockUserOutlookAPI{}
	realUserService := *NewUserService(mockUserRepo, mockOutlook)

	svc := NewCommentService(mockRepo, realUserService, mockImageStorage, mockFileUtils, "bucket123")

	req := &domain.CreateCommentReq{
		SessionID: "u1",
		ImageData: []*multipart.FileHeader{{Filename: "file1.png"}},
	}

	_, err := svc.CreateComment(context.Background(), req)
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected 'user not found', got %v", err)
	}
}

func TestLoadComments_Success(t *testing.T) {
	expected := []*domain.Comment{
		{Content: "c1"},
		{Content: "c2"},
	}
	mockRepo := &MockCommentRepo{comments: expected}

	svc := NewCommentService(mockRepo, UserService{}, nil, nil, "")

	got, err := svc.LoadComments(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %+v, got %+v", expected, got)
	}
}
