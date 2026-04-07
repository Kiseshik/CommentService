package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type PostRepository struct {
	mu    sync.RWMutex
	store map[string]*domain.Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{ //redis
		store: make(map[string]*domain.Post),
	}
}

func (r *PostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return nil, domain.ErrPostNotFound
	}
	return p, nil
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.store[post.ID]; exists {
		return domain.ErrPostAlreadyExists
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
		return domain.ErrPostNotFound
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
		return domain.ErrPostNotFound
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

func (r *PostRepository) List(ctx context.Context) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	posts := make([]*domain.Post, 0, len(r.store))
	for _, p := range r.store {
		posts = append(posts, p)
	}
	sortByCreatedAt(posts)
	return posts, nil
}
