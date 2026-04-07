package domain

import (
	"errors"
	"time"
)

const (
	MaxPostTitleLength   = 200
	MaxPostContentLength = 10000
)

var (
	ErrPostNotFound      = errors.New("post not found")
	ErrPostAlreadyExists = errors.New("post already exists")
)

type Post struct {
	ID              string    `db:"id"`
	Title           string    `db:"title"`
	Content         string    `db:"content"`
	AuthorID        string    `db:"author_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	CommentsEnabled bool      `db:"comments_enabled"`
}

func (p *Post) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Post) GetID() string {
	return p.ID
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("post title is required")
	}
	if len(p.Title) > MaxPostTitleLength {
		return errors.New("post title exceeds maximum length")
	}
	if p.Content == "" {
		return errors.New("post content is required")
	}
	if len(p.Content) > MaxPostContentLength {
		return errors.New("post content exceeds maximum length")
	}
	if p.AuthorID == "" {
		return errors.New("post author_id is required")
	}
	return nil
}
