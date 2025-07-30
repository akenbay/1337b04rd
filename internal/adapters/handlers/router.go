package handlers

import (
	"1337b04rd/internal/services"
	"net/http"
)

func NewRouter(userService services.UserService, postService services.PostService, commentService services.CommentService) *http.ServeMux {
	mux := http.NewServeMux()
	userHandler := newUserHandlers(userService)
	postHandler := newPostHandlers(postService)
	commentHandler := newCommentHandlers(commentService)

	mux.HandleFunc("GET /session/me", userHandler.getSessionMe)
	mux.HandleFunc("GET /threads", postHandler.getActivePostsApi)
	mux.HandleFunc("GET /threads/archive", postHandler.getArchivedPostsApi)
	mux.HandleFunc("GET /threads/view/", postHandler.getPostApi)
	mux.HandleFunc("POST /threads", postHandler.createPostAPI)
	mux.HandleFunc("POST /threads/comment", commentHandler.createCommentAPI)
	mux.HandleFunc("GET /threads/comment", commentHandler.loadCommentsApi)

	return mux
}
