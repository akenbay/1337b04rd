package post

import "errors"

var (
	ErrPostTitleRequired   = errors.New("post title is required")
	ErrPostContentRequired = errors.New("post content is required")
	ErrPostNotFound        = errors.New("post not found")
	ErrImageTooLarge       = errors.New("image size exceeds limit")
)
