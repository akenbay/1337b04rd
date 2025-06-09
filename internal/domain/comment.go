package domain

import "time"

type Comment struct {
	ID        string  // UUID
	PostID    string  // Parent post UUID
	ParentID  *string // Nullable (for nested comments)
	User      User    // Embedded or reference SessionID
	Content   string
	CreatedAt time.Time
}
