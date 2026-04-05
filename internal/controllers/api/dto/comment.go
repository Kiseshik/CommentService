package dto

import "time"

type ListCommentsRequest struct {
	PostID   string  `json:"postId"`
	ParentID *string `json:"parentId,omitempty"`
	Limit    int     `json:"limit"`
	Cursor   string  `json:"cursor"`
}

type CreateCommentRequest struct {
	PostID   string  `json:"postId"`
	ParentID *string `json:"parentId,omitempty"`
	AuthorID string  `json:"authorId"`
	Content  string  `json:"content"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	ParentID  *string   `json:"parentId,omitempty"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommentListResponse struct {
	Comments    []CommentResponse `json:"comments"`
	HasNextPage bool              `json:"hasNextPage"`
	NextCursor  string            `json:"nextCursor"`
}
