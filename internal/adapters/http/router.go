package http

import (
	"1337b04rd/internal/services"
	"net/http"
)

func NewRouter(userService services.UserService, postService services.PostService) *http.ServeMux {
	mux := http.NewServeMux()
	userHandler := newUserHandlers(userService)
	postHandler := newPostHandlers(postService)

	mux.HandleFunc("GET /session/me")
}
