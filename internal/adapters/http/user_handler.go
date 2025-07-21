package http

import (
	"1337b04rd/internal/services"
	"net/http"
)

type UserHandlers struct {
	userService services.UserService
}

func NewUserHandlers(userService services.UserService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

func (u *UserHandlers) GetSessionMe(w http.ResponseWriter, r *http.Request) {
	ID, err := getSessionID(r)
	if err != nil {
		respondError(w, "Failed to retrieve the user session", http.StatusNotFound)
	}

	user, err := u.userService.FindUserByID(r.Context(), ID)
	if err != nil {
		respondError(w, "Failed to retrieve the user session", http.StatusNotFound)
	}
	respondJSON(w, user, http.StatusOK)
}
