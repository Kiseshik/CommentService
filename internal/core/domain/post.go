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
	ErrPostNotFound            = errors.New("post not found")
	ErrPostAlreadyExists       = errors.New("post already exists")
	ErrParentFromDifferentPost = errors.New("parent comment belongs to different post")
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
