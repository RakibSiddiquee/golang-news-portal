package models

import "time"

type Post struct {
	ID           int       `json:"id,omitempty"`
	Title        string    `json:"title"`
	Url          string    `json:"url"`
	CreateAt     time.Time `json:"create_at"`
	UserID       int       `json:"user_id"`
	Votes        int       `json:"votes,omitempty"`
	UserName     string    `json:"user_name,omitempty"`
	CommentCount int       `json:"comment_count,omitempty"`
	TotalRecords int       `json:"total_records,omitempty"`
}
