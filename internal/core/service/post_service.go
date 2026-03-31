package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/google/uuid"
)

type PostService struct {
	postRepo port.PostRepository
}

func NewPostService(postRepo port.PostRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
	}
}

func (s *PostService) CreatePost(ctx context.Context, title, content, authorID string, commentsEnabled bool) (*domain.Post, error) {
	post := &domain.Post{
		ID:              uuid.New().String(),
		Title:           title,
		Content:         content,
		AuthorID:        authorID,
		CommentsEnabled: commentsEnabled,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := post.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	return post, nil
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

func (s *PostService) ToggleComments(ctx context.Context, postID string, enabled bool) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	post.CommentsEnabled = enabled
	post.UpdatedAt = time.Now()
	if err := s.postRepo.Update(ctx, post); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (s *PostService) Exists(ctx context.Context, id string) (bool, error) {
	return s.postRepo.Exists(ctx, id)
}
