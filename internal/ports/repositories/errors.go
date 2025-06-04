package repositories

import "errors"

var (
	ErrUserNotFound = errors.New("User was not found")
	ErrPostNotFound = errors.New("Post was not found")
)
