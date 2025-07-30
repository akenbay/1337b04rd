package handlers

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"errors"
	"log/slog"
	"net/http"
)

type PostHandlers struct {
	postService services.PostService
}

func newPostHandlers(postService services.PostService) *PostHandlers {
	return &PostHandlers{
		postService: postService,
	}
}

func (h *PostHandlers) createPostAPI(w http.ResponseWriter, r *http.Request) {
	slog.Info("Creating post handler:")

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		respondError(w, r, "Failed to parse files, the size is too big", http.StatusBadRequest)
		return
	}

	sessionID, err := getSessionID(r)
	if err != nil {
		respondError(w, r, "Failed to get session id from cookies", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")     // Works after ParseMultipartForm
	content := r.FormValue("content") // Same here

	files := r.MultipartForm.File["images"]

	post, err := h.postService.CreatePost(r.Context(), &domain.CreatePostReq{
		Title:     title,
		Content:   content,
		ImageData: files,
		SessionID: sessionID,
	})
	if err != nil {
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, r, post, http.StatusCreated)
}

func (h *PostHandlers) getPostApi(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	postID := r.URL.Path[len("/threads/"):] // Gets everything after "/posts/"
	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, "Post not found", http.StatusNotFound)
			return
		}
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, r, post, http.StatusOK)
	return
}

func (h *PostHandlers) getActivePostsApi(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			respondError(w, r, "Internal server error", http.StatusInternalServerError)
			return
		}
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, r, posts, http.StatusOK)
	return
}

func (h *PostHandlers) getArchivedPostsApi(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetArchivedPosts(r.Context())
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			respondError(w, r, "Internal server error", http.StatusInternalServerError)
			return
		}
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, r, posts, http.StatusOK)
	return
}
