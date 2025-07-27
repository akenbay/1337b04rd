package handlers

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
)

type CommentHandlers struct {
	commentService services.CommentService
}

func newCommentHandlers(commentService services.CommentService) *CommentHandlers {
	return &CommentHandlers{
		commentService: commentService,
	}
}

func (h *CommentHandlers) createCommentAPI(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content string `json:"content"`
		Image   string `json:"image"` // base64 encoded
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, r, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Process base64 image if provided
	var imageData []byte
	if req.Image != "" {
		var err error
		imageData, err = base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			respondError(w, r, "Invalid image encoding", http.StatusBadRequest)
			return
		}
	}

	sessionID, err := getSessionID(r)
	if err != nil {
		respondError(w, r, "Failed to get session id from cookies", http.StatusBadRequest)
		return
	}

	comment, err := h.commentService.CreateComment(r.Context(), &domain.CreateCommentReq{
		Content:   req.Content,
		ImageData: imageData,
		SessionID: sessionID,
	})
	if err != nil {
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, r, comment, http.StatusCreated)
}

func (h *CommentHandlers) loadCommentsApi(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	postID := r.PathValue("threadId")
	comments, err := h.commentService.LoadComments(r.Context(), postID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, "Post not found", http.StatusNotFound)
			return
		}
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, r, comments, http.StatusOK)
	return
}

// func (h *PostHandlers) getActivePostsApi(w http.ResponseWriter, r *http.Request) {
// 	posts, err := h.postService.GetActivePosts(r.Context())
// 	if err != nil {
// 		if r.Header.Get("Accept") == "application/json" {
// 			respondError(w, r, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
// 		respondError(w, r, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	respondJSON(w, r, posts, http.StatusOK)
// 	return
// }

// func (h *PostHandlers) getArchivedPostsApi(w http.ResponseWriter, r *http.Request) {
// 	posts, err := h.postService.GetArchivedPosts(r.Context())
// 	if err != nil {
// 		if r.Header.Get("Accept") == "application/json" {
// 			respondError(w, r, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
// 		respondError(w, r, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	respondJSON(w, r, posts, http.StatusOK)
// 	return
// }
