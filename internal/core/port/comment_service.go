package port

import (
	"context"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type CommentService interface {
	CreateComment(ctx context.Context, params *CreateCommentParams) (*domain.Comment, error)
	GetCommentByID(ctx context.Context, id string) (*domain.Comment, error)
	ListComments(ctx context.Context, params *ListCommentParams) (*CommentListResult, error)
}

type CreateCommentParams struct {
	PostID   string
	ParentID *string
	AuthorID string
	Content  string
}

type ListCommentParams struct {
	PostID   string
	ParentID *string
	Cursor   string
	Limit    int
}
