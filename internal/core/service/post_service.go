package service

import (
	"context"
	"fmt"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
)

type PostService struct {
	postRepo port.PostRepository
}

func NewPostService(postRepo port.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) CreatePost(ctx context.Context, params *port.CreatePostParams) (*domain.Post, error) {
	tempPost := &domain.Post{
		Title:           params.Title,
		Content:         params.Content,
		AuthorID:        params.AuthorID,
		CommentsEnabled: params.CommentsEnabled,
	}
	if err := tempPost.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}
	return s.postRepo.Create(ctx, params)
}

func (s *PostService) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

func (s *PostService) ListPosts(ctx context.Context) ([]*domain.Post, error) {
	posts, err := s.postRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	return posts, nil
}

func (s *PostService) UpdatePost(ctx context.Context, params *port.UpdatePostParams) (*domain.Post, error) {
	post, err := s.postRepo.Update(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	return post, nil
}

func (s *PostService) ToggleComments(ctx context.Context, postID string) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	newEnabled := !post.CommentsEnabled
	_, err = s.postRepo.Update(ctx, &port.UpdatePostParams{
		ID:              postID,
		CommentsEnabled: &newEnabled,
	})
	if err != nil {
		return fmt.Errorf("failed to toggle comments: %w", err)
	}
	return nil
}

func (s *PostService) Exists(ctx context.Context, id string) (bool, error) {
	return s.postRepo.Exists(ctx, id)
}
