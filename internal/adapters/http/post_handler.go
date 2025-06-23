package http

import (
	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type PostHandlers struct {
	postService    services.PostService
	imageValidator ImageValidator
	templates      *template.Template
}

type ImageValidator interface {
	Validate(image []byte) error
	AllowedTypes() []string
}

func NewPostHandlers(postService services.PostService, validator ImageValidator, templateDir string) *PostHandlers {
	templates := template.Must(template.ParseGlob(filepath.Join(templateDir, "*.html")))
	return &PostHandlers{
		postService:    postService,
		imageValidator: validator,
		templates:      templates,
	}
}

// Render HTML Responses

func (h *PostHandlers) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("Template rendering error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *PostHandlers) renderError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	h.renderTemplate(w, "error.html", map[string]interface{}{
		"Error":  message,
		"Status": status,
	})
}

// JSON Response Helpers

func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

func respondError(w http.ResponseWriter, message string, status int) {
	respondJSON(w, map[string]string{"error": message}, status)
}

// Handlers

func (h *PostHandlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Handle both form and JSON input
	if r.Header.Get("Content-Type") == "application/json" {
		h.createPostAPI(w, r)
		return
	}

	// HTML form handling
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		h.renderError(w, "File too large", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	if len(title) < 5 || len(content) < 1 {
		h.renderError(w, "Title (min 5 chars) and content required", http.StatusBadRequest)
		return
	}

	// Process image upload
	var imageData []byte
	if file, header, err := r.FormFile("image"); err == nil {
		defer file.Close()

		// Basic checks
		if header.Size > 10<<20 { // 10MB
			h.renderError(w, "Image too large", http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		switch contentType {
		case "image/jpeg", "image/png":
			// Allowed types
		default:
			h.renderError(w, "Only JPEG/PNG allowed", http.StatusBadRequest)
			return
		}

		imageData, err = io.ReadAll(file)
		if err != nil {
			h.renderError(w, "Failed to read image", http.StatusBadRequest)
			return
		}

		if err := h.imageValidator.Validate(imageData); err != nil {
			h.renderError(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		log.Printf("Post creation failed: %v", err)
		h.renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := h.postService.CreatePost(r.Context(), &domain.CreatePostReq{
		Title:     title,
		Content:   content,
		ImageData: imageData,
	})
	if err != nil {
		log.Printf("Post creation failed: %v", err)
		h.renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Redirect to the new post
	http.Redirect(w, r, "/posts/"+id, http.StatusSeeOther)
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

func (h *PostHandlers) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.renderError(w, "Post not found", http.StatusNotFound)
			return
		}
		h.renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check Accept header for response format
	if r.Header.Get("Accept") == "application/json" {
		respondJSON(w, post, http.StatusOK)
		return
	}

	h.renderTemplate(w, "post.html", map[string]interface{}{
		"Post": post,
		// Additional template data
	})
}

func (h *PostHandlers) ListActivePosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			respondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		h.renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondJSON(w, posts, http.StatusOK)
		return
	}

	h.renderTemplate(w, "catalog.html", map[string]interface{}{
		"Posts": posts,
	})
}

func (h *PostHandlers) ListArchivedPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetArchivedPosts(r.Context())
	if err != nil {
		if r.Header.Get("Accept") == "application/json" {
			respondError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		h.renderError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/json" {
		respondJSON(w, posts, http.StatusOK)
		return
	}

	h.renderTemplate(w, "catalog.html", map[string]interface{}{
		"Posts": posts,
	})
}
