package port

import (
	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type PubSub interface {
	Subscribe(postID string) (<-chan *domain.Comment, func())
	Publish(postID string, comment *domain.Comment)
}
