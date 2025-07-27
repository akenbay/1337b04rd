package handlers

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"encoding/base64"
	"encoding/json"
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
	slog.Info("Creating post API:")

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Image   string `json:"image"` // base64 encoded
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("Successfully decoded request body")

	// Process base64 image if provided
	var imageData []byte
	if req.Image != "" {
		var err error
		imageData, err = base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			respondError(w, r, "Invalid image encoding", http.StatusBadRequest)
			return
		}
		slog.Info("Decoded imagedata successfully")
	}

	sessionID, err := getSessionID(r)
	if err != nil {
		respondError(w, r, "Failed to get session id from cookies", http.StatusBadRequest)
		return
	}

	post, err := h.postService.CreatePost(r.Context(), &domain.CreatePostReq{
		Title:     req.Title,
		Content:   req.Content,
		ImageData: imageData,
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
