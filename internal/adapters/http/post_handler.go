package http

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

type PostHandlers struct {
	postService    services.PostService
	imageValidator ImageValidator
}

type ImageValidator interface {
	Validate(image []byte) error
	AllowedTypes() []string
}

func NewPostHandlers(postService services.PostService, validator ImageValidator, templateDir string) *PostHandlers {
	return &PostHandlers{
		postService:    postService,
		imageValidator: validator,
	}
}

func (h *PostHandlers) createPostAPI(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Image   string `json:"image"` // base64 encoded
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process base64 image if provided
	var imageData []byte
	if req.Image != "" {
		var err error
		imageData, err = base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			respondError(w, "Invalid image encoding", http.StatusBadRequest)
			return
		}

		if err := h.imageValidator.Validate(imageData); err != nil {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	post, err := h.postService.CreatePost(r.Context(), &domain.CreatePostReq{
		Title:     req.Title,
		Content:   req.Content,
		ImageData: imageData,
	})
	if err != nil {
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, post, http.StatusCreated)
}

func (h *PostHandlers) GetPostApi(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, "Post not found", http.StatusNotFound)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, post, http.StatusOK)
	return
}

func (h *PostHandlers) GetActivePostsApi(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			respondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, posts, http.StatusOK)
	return
}

func (h *PostHandlers) GetArchivedPostsApi(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetArchivedPosts(r.Context())
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			respondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		respondError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, posts, http.StatusOK)
	return
}
