package port

import (
	"context"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type PostRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Post, error)
	Create(ctx context.Context, params *CreatePostParams) (*domain.Post, error)
	Update(ctx context.Context, params *UpdatePostParams) (*domain.Post, error)
	List(ctx context.Context) ([]*domain.Post, error)
	Exists(ctx context.Context, id string) (bool, error)
}
