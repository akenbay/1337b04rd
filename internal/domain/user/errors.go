package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidSession    = errors.New("invalid session")
	ErrInvalidName       = errors.New("name must be 1-50 characters")
	ErrAvatarUnavailable = errors.New("could not fetch avatar")
)
