package port

import (
	"context"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type CommentRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Comment, error)
	Create(ctx context.Context, params *CreateCommentParams) (*domain.Comment, error)
	Update(ctx context.Context, params *UpdateCommentParams) (*domain.Comment, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, params CommentListParams) (*CommentListResult, error)
	CountByPost(ctx context.Context, postID string) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
}
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

type UpdateCommentParams struct {
	ID      string
	Content *string
}
