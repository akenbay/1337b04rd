package domain

import "time"

type Post struct {
	ID        string
	Title     string
	Content   string
	ImageURL  string
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}
