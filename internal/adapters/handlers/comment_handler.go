package handlers

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"errors"
	"log/slog"
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
	slog.Info("API creating comment:")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error("Could not parse images attached to comment:", "error", err)
		respondError(w, r, "Failed to parse attached files (maximum size of files 10 MB)", http.StatusBadRequest)
	}

	slog.Info("Parsed multipart form")

	content := r.FormValue("content")
	postID := r.FormValue("thread_id")
	files := r.MultipartForm.File["images"]

	sessionID, err := getSessionID(r)
	if err != nil {
		respondError(w, r, "Failed to get session id from cookies", http.StatusBadRequest)
		return
	}

	slog.Info("Got sessionID")

	comment, err := h.commentService.CreateComment(r.Context(), &domain.CreateCommentReq{
		Content:   content,
		ImageData: files,
		SessionID: sessionID,
		PostID:    postID,
	})
	if err != nil {
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Created comment")

	respondJSON(w, r, comment, http.StatusCreated)
}

func (h *CommentHandlers) loadCommentsApi(w http.ResponseWriter, r *http.Request) {
	slog.Info("API load comments:")

	// Extract the ID from the URL path
	postID := r.URL.Query().Get("thread_id")

	slog.Info("Got id from URL:", "id", postID)

	comments, err := h.commentService.LoadComments(r.Context(), postID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, "Post not found", http.StatusNotFound)
			return
		}
		slog.Error("Error when loading comments:", "error", err)
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Loaded comments")

	respondJSON(w, r, comments, http.StatusOK)
	return
}
