package dto

import "time"

type CreatePostRequest struct {
	Title           string `json:"title"`
	Content         string `json:"content"`
	AuthorID        string `json:"authorId"`
	CommentsEnabled bool   `json:"commentsEnabled"`
}

type UpdatePostRequest struct {
	Title           *string `json:"title,omitempty"`
	Content         *string `json:"content,omitempty"`
	CommentsEnabled *bool   `json:"commentsEnabled,omitempty"`
}

type PostResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	AuthorID        string    `json:"authorId"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
