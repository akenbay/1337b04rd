package handlers

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// JSON Response Helpers

func respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

func respondError(w http.ResponseWriter, r *http.Request, message string, status int) {
	respondJSON(w, r, map[string]string{"error": message}, status)
}

// Retrieve session ID from cookies

func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Error("Failed to get session_id from cookies", "error", err)
		return "", err
	}
	return cookie.Value, nil
}

// Set session ID to the cookies

func setSessionID(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour * 7),
	}
	http.SetCookie(w, cookie)
}
