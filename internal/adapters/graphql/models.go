package graphql

import "github.com/Kiseshik/CommentService.git/internal/core/domain"

type CommentListResult struct {
	Comments    []*domain.Comment `json:"comments"`
	HasNextPage bool              `json:"hasNextPage"`
	NextCursor  *string           `json:"nextCursor"`
}
