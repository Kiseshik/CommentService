package domain

import (
	"errors"
	"time"
)

const (
	MaxCommentLength = 2000
)

var (
	ErrCommentNotFound      = errors.New("comment not found")
	ErrEmptyComment         = errors.New("comment cannot be empty")
	ErrCommentAlreadyExists = errors.New("comment already exists")
	ErrCommentTooLong       = errors.New("comment exceeds maximum length of 2000 characters")
	ErrCommentsDisabled     = errors.New("comments are disabled for this post")
	ErrParentNotFound       = errors.New("parent comment not found")
	ErrMaxDepthExceeded     = errors.New("max depth exceeded")
)

type Comment struct {
	ID        string    `db:"id"`
	PostID    string    `db:"post_id"`
	ParentID  *string   `db:"parent_id"`
	Content   string    `db:"content"`
	AuthorID  string    `db:"author_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Comment) GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c *Comment) GetID() string {
	return c.ID
}

func (c *Comment) Validate() error {
	if c.PostID == "" {
		return errors.New("comment post_id is required")
	}
	if c.Content == "" {
		return errors.New("comment content is required")
	}
	if len(c.Content) > MaxCommentLength {
		return errors.New("comment content exceeds maximum length")
	}
	if c.AuthorID == "" {
		return errors.New("comment author_id is required")
	}
	return nil
}
