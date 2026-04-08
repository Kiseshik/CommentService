package port

import (
	"context"
	"errors"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type PostService interface {
	CreatePost(ctx context.Context, params *CreatePostParams) (*domain.Post, error)
	GetPostByID(ctx context.Context, id string) (*domain.Post, error)
	ListPosts(ctx context.Context) ([]*domain.Post, error)
	UpdatePost(ctx context.Context, params *UpdatePostParams) (*domain.Post, error)
	ToggleComments(ctx context.Context, postID string) error
	Exists(ctx context.Context, id string) (bool, error)
}

type CreatePostParams struct {
	Title           string
	Content         string
	AuthorID        string
	CommentsEnabled bool
}

func (p *CreatePostParams) Validate() error {
	if p.Title == "" {
		return errors.New("post title is required")
	}
	if len(p.Title) > domain.MaxPostTitleLength {
		return errors.New("post title exceeds maximum length")
	}
	if p.Content == "" {
		return errors.New("post content is required")
	}
	if len(p.Content) > domain.MaxPostContentLength {
		return errors.New("post content exceeds maximum length")
	}
	if p.AuthorID == "" {
		return errors.New("post author_id is required")
	}
	return nil
}

type UpdatePostParams struct {
	ID              string
	Title           *string
	Content         *string
	CommentsEnabled *bool
}
