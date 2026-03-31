package port

import (
	"context"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type PostListParams struct {
	AuthorID        string
	Title           string
	CommentsEnabled *bool
	Cursor          string
	Limit           int
}
type PostListResult struct {
	Posts       []*domain.Post
	NextCursor  string
	HasNextPage bool
}

type PostRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	List(ctx context.Context, params PostListParams) (*PostListResult, error)
	Exists(ctx context.Context, id string) (bool, error)
}
