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

// Get the current user session, if it does not exists creating and retrieving the newly created one

func (u *UserHandlers) GetSessionMe(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionID(r)
	if err != nil {
		sessionID, err = u.createUser(r)
		if err != nil {
			respondError(w, "Failed to create new user session", http.StatusNotFound)
			return
		}
	}

	user, err := u.userService.FindUserByID(r.Context(), sessionID)
	if err != nil {
		sessionID, err = u.createUser(r)
		if err != nil {
			respondError(w, "Failed to create new user session", http.StatusNotFound)
			return
		}
		setSessionID(w, sessionID)
		user, err := u.userService.FindUserByID(r.Context(), sessionID)
		if err != nil {
			respondError(w, "Failed to retrieve the user session", http.StatusNotFound)
			return
		}
		respondJSON(w, user, http.StatusOK)
		return
	}
	respondJSON(w, user, http.StatusOK)
}

// Create new user and retrieve its ID in the database

func (u *UserHandlers) createUser(r *http.Request) (string, error) {
	sessionID, err := u.userService.CreateUserAndGetID(r.Context())
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
