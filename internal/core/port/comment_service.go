package port

import (
	"context"
	"errors"

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

func (p *CreateCommentParams) Validate() error {
	if p.PostID == "" {
		return errors.New("comment post_id is required")
	}
	if p.Content == "" {
		return domain.ErrEmptyComment
	}
	if len(p.Content) > domain.MaxCommentLength {
		return domain.ErrCommentTooLong
	}
	if p.AuthorID == "" {
		return errors.New("comment author_id is required")
	}
	return nil
}
