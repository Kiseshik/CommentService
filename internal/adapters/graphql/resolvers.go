package graphql

import (
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/Kiseshik/CommentService.git/internal/core/service"
)

type Resolver struct {
	postService    *service.PostService
	commentService *service.CommentService
	pubsub         port.PubSub
}
type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

func NewResolver(
	postService *service.PostService,
	commentService *service.CommentService,
	pubsub port.PubSub,
) *Resolver {
	return &Resolver{
		postService:    postService,
		commentService: commentService,
		pubsub:         pubsub,
	}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}
