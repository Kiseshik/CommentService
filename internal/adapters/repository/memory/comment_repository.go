package memory

import (
	"context"
	"sync"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
)

type CommentRepository struct {
	mu    sync.RWMutex
	store map[string]*domain.Comment
	//postIndex  map[string][]string
	//parentIndex map[string][]string
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		store: make(map[string]*domain.Comment), //сюда редис
	}
}

func (r *CommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.store[id]
	if !ok {
		return nil, domain.ErrCommentNotFound
	}
	return c, nil
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.store[comment.ID]; exists {
		return domain.ErrCommentAlreadyExists
	}
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	r.store[comment.ID] = comment
	return nil
}

func (r *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existingCom, exists := r.store[comment.ID]
	if !exists {
		return domain.ErrCommentNotFound
	}
	comment.CreatedAt = existingCom.CreatedAt
	comment.UpdatedAt = time.Now()
	r.store[comment.ID] = comment
	return nil
}

func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.store[id]; !exists {
		return domain.ErrCommentNotFound
	}
	delete(r.store, id)
	return nil
}

func (r *CommentRepository) Exists(ctx context.Context, id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.store[id]
	return exists, nil
}

func (r *CommentRepository) CountByPost(ctx context.Context, postID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, c := range r.store {
		if c.PostID == postID {
			count++
		}
	}
	return count, nil
}

func (r *CommentRepository) List(ctx context.Context, params port.CommentListParams) (*port.CommentListResult, error) {
	r.mu.RLock() //вероятно ботлнек для хабров и прочих хайлоадов, надо бы чего-то придумывать, инексация?
	defer r.mu.RUnlock()

	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}

	count, err := r.CountByPost(ctx, params.PostID) // с расчетом на индексацию, пока так
	if err != nil {
		return nil, err
	}
	filtered := make([]*domain.Comment, 0, count)
	for _, c := range r.store {
		if c.PostID != params.PostID {
			continue
		}
		if params.ParentID != nil && (c.ParentID == nil || *c.ParentID != *params.ParentID) {
			continue
		}
		if params.ParentID == nil && c.ParentID != nil {
			continue
		}
		filtered = append(filtered, c)
	}
	sortByCreatedAt(filtered)
	return paginateComments(filtered, params)
}

func paginateComments(comments []*domain.Comment, params port.CommentListParams) (*port.CommentListResult, error) {
	start := 0
	if params.Cursor != "" {
		found := false
		for i, c := range comments {
			if c.ID == params.Cursor {
				start = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, domain.ErrInvalidCursor
		}
	}
	end := start + params.Limit
	if end > len(comments) {
		end = len(comments)
	}
	res := &port.CommentListResult{
		Comments:    comments[start:end],
		HasNextPage: end < len(comments),
	}
	if res.HasNextPage {
		res.NextCursor = comments[end-1].ID
	}
	return res, nil
}
