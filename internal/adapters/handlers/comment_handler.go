package handlers

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"1337b04rd/pkg/logger"
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
	logger.Info("API creating comment:")

	createReq := domain.CreateCommentReq{}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error("Could not parse images attached to comment:", "error", err)
		respondError(w, r, "Failed to parse attached files (maximum size of files 10 MB)", http.StatusBadRequest)
	}

	logger.Info("Parsed multipart form")

	createReq.Content = r.FormValue("content")
	createReq.PostID = r.FormValue("thread_id")

	if parentID := r.FormValue("parent_id"); parentID != "" {
		createReq.ParentID = &parentID
		logger.Info("Found reply comment:", "parent id", parentID)
	}

	createReq.ImageData = r.MultipartForm.File["images"]

	createReq.SessionID, err = getSessionID(r)
	if err != nil {
		respondError(w, r, "Failed to get session id from cookies", http.StatusBadRequest)
		return
	}

	logger.Info("Got sessionID")

	comment, err := h.commentService.CreateComment(r.Context(), &createReq)
	if err != nil {
		respondError(w, r, err.Error(), 200)
		return
	}

	logger.Info("Created comment")

	respondJSON(w, r, comment, http.StatusCreated)
}

func (h *CommentHandlers) loadCommentsApi(w http.ResponseWriter, r *http.Request) {
	logger.Info("API load comments:")

	// Extract the ID from the URL path
	postID := r.URL.Query().Get("thread_id")

	logger.Info("Got id from URL:", "id", postID)

	comments, err := h.commentService.LoadComments(r.Context(), postID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, r, "Post not found", http.StatusNotFound)
			return
		}
		logger.Error("Error when loading comments:", "error", err)
		respondError(w, r, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info("Loaded comments")

	respondJSON(w, r, comments, http.StatusOK)
	return
}
