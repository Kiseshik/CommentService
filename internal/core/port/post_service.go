package port

import (
	"context"

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

type UpdatePostParams struct {
	ID              string
	Title           *string
	Content         *string
	CommentsEnabled *bool
}
