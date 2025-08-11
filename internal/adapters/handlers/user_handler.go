package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/services"
)

type UserHandlers struct {
	userService services.UserService
}

func newUserHandlers(userService services.UserService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

// Get the current user session, if it does not exists creating and retrieving the newly created one

func (u *UserHandlers) getSessionMe(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionID(r)
	if err != nil {
		sessionID, err = u.createUser(r)
		if err != nil {
			respondError(w, r, "Failed to create new user session", http.StatusNotFound)
			return
		}
		slog.Info("Successfuly created user")
		setSessionID(w, sessionID)
		slog.Info("Successfuly set session id")
	}

	user, err := u.userService.FindUserByID(r.Context(), sessionID)
	if err != nil {
		slog.Info("Error when finding user")
		sessionID, err = u.createUser(r)
		if err != nil {
			respondError(w, r, "Failed to create new user session", http.StatusNotFound)
			return
		}
		slog.Info("Successfuly created user")
		setSessionID(w, sessionID)
		user, err := u.userService.FindUserByID(r.Context(), sessionID)
		if err != nil {
			respondError(w, r, "Failed to retrieve the user session", http.StatusNotFound)
			return
		}
		respondJSON(w, r, user, http.StatusOK)
		return
	}
	slog.Info("Found user by session id")
	respondJSON(w, r, user, http.StatusOK)
}

// Create new user and retrieve its ID from the database

func (u *UserHandlers) createUser(r *http.Request) (string, error) {
	sessionID, err := u.userService.CreateUserAndGetID(r.Context())
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (u *UserHandlers) changeUsername(w http.ResponseWriter, r *http.Request) {
	slog.Info("Change username handler:")

	sessionID, err := getSessionID(r)

	var req domain.NameRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("Error when decoding new username:", "error", err)
		respondError(w, r, "Invalid username", http.StatusBadRequest)
		return
	}
	newUsername := req.DisplayName

	err = u.userService.ChangeUsername(r.Context(), sessionID, newUsername)
	if err != nil {
		slog.Error("Error when changing username:", "error", err)
	}
}
