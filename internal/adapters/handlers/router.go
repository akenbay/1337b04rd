package handlers

import (
	"1337b04rd/internal/services"
	"net/http"
)

func NewRouter(userService services.UserService, postService services.PostService) *http.ServeMux {
	mux := http.NewServeMux()
	userHandler := newUserHandlers(userService)
	postHandler := newPostHandlers(postService)

	mux.HandleFunc("GET /session/me", userHandler.getSessionMe)
	mux.HandleFunc("GET /threads", postHandler.getActivePostsApi)
	mux.HandleFunc("GET /threads/archive", postHandler.getArchivedPostsApi)
	mux.HandleFunc("GET /threads/view/", postHandler.getPostApi)
	mux.HandleFunc("POST /threads", postHandler.createPostAPI)

	return mux
}
