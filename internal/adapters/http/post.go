package http

import (
	"1337b04rd/internal/services"
	"net/http"
)

type PostHandlers struct {
	postService services.PostService
}

func (h *PostHandlers) Create(w http.ResponseWriter, r *http.Request) {
}
