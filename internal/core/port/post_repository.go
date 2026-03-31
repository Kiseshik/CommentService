package port

import (
	"context"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type PostListParams struct {
	Cursor string
	Limit  int
}

type PostRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	List(ctx context.Context, params PostListParams) ([]*domain.Post, string, error)
	Exists(ctx context.Context, id string) (bool, error)
}
