package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/google/uuid"
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

func (r *PostRepository) Create(ctx context.Context, params *port.CreatePostParams) (*domain.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	post := &domain.Post{
		ID:              uuid.New().String(),
		Title:           params.Title,
		Content:         params.Content,
		AuthorID:        params.AuthorID,
		CommentsEnabled: params.CommentsEnabled,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	r.store[post.ID] = post
	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, params *port.UpdatePostParams) (*domain.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.store[params.ID]
	if !ok {
		return nil, domain.ErrPostNotFound
	}

	if params.Title != nil {
		existing.Title = *params.Title
	}
	if params.Content != nil {
		existing.Content = *params.Content
	}
	if params.CommentsEnabled != nil {
		existing.CommentsEnabled = *params.CommentsEnabled
	}
	existing.UpdatedAt = time.Now() //кажется это чище чем прописывать апдейтит эт в бизнес слое, разве нет?

	r.store[existing.ID] = existing
	return existing, nil
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
