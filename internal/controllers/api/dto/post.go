package dto

import "time"

type CreatePostRequest struct {
	Title           string `json:"title"`
	Content         string `json:"content"`
	AuthorID        string `json:"authorId"`
	CommentsEnabled bool   `json:"commentsEnabled"`
}

type GetPostRequest struct {
	ID string `json:"id"`
}

type UpdatePostRequest struct {
	ID              string  `json:"id"`
	Title           *string `json:"title,omitempty"`
	Content         *string `json:"content,omitempty"`
	CommentsEnabled *bool   `json:"commentsEnabled,omitempty"`
}

type ToggleCommentsRequest struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
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
