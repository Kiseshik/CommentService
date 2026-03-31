package port

import (
	"context"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type CommentListParams struct {
	PostID   string
	ParentID *string
	Cursor   string
	Limit    int
}

type CommentListResult struct {
	Comments    []*domain.Comment
	NextCursor  string
	HasNextPage bool
}

type CommentRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Comment, error)
	Create(ctx context.Context, comment *domain.Comment) error
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params CommentListParams) (*CommentListResult, error)
	ListAllByPost(ctx context.Context, postID string) ([]*domain.Comment, error)
	CountByPost(ctx context.Context, postID string) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
}
