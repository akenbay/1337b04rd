package domain

import "time"

type Comment struct {
	ID        string
	PostID    string
	ParentID  string
	Content   string
	User      User
	CreatedAt time.Time
}
