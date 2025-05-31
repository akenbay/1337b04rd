package post

import (
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments,omitempty"`
}

func (p *Post) IsExpired() bool {
	expiryDuration := 10 * time.Minute

	if len(p.Comments) > 0 {
		lastActivity := p.UpdatedAt
		for _, c := range p.Comments {
			if c.CreatedAt.After(lastActivity) {
				lastActivity = c.CreatedAt
			}
		}
		expiryDuration = 15 * time.Minute
		return time.Now().After(lastActivity.Add(expiryDuration))
	}

	return time.Now().After(p.CreatedAt.Add(expiryDuration))
}
