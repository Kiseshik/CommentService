package memory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/Kiseshik/CommentService.git/internal/utils"
)

type PostRepository struct {
	mu    sync.RWMutex
	store map[string]*domain.Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		store: make(map[string]*domain.Post),
	}
}

func (r *PostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return p, nil
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.store[post.ID]; exists {
		return errors.New("post already exists")
	}
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now
	r.store[post.ID] = post
	return nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingPost, exists := r.store[post.ID]
	if !exists {
		return errors.New("post not found")
	}
	post.CreatedAt = existingPost.CreatedAt
	post.UpdatedAt = time.Now()
	r.store[post.ID] = post
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.store[id]; !exists {
		return errors.New("post not found")
	}
	delete(r.store, id)
	return nil
}

func (r *PostRepository) Exists(ctx context.Context, id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.store[id]
	return exists, nil
}

func (r *PostRepository) List(ctx context.Context, params port.PostListParams) (*port.PostListResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	filtered := make([]*domain.Post, 0, len(r.store))
	for _, p := range r.store {
		if params.AuthorID != "" && p.AuthorID != params.AuthorID {
			continue
		}
		if params.Title != "" && !strings.Contains(p.Title, params.Title) {
			continue
		}
		if params.CommentsEnabled != nil && p.CommentsEnabled != *params.CommentsEnabled {
			continue
		}
		filtered = append(filtered, p)
	}
	utils.SortPosts(filtered)
	return utils.PaginatePosts(filtered, params)
}
