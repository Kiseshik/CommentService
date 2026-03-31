package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/google/uuid"
)

const (
	MaxCommentDepth = 100
)

type CommentService struct {
	commentRepo port.CommentRepository
	postRepo    port.PostRepository
}

func NewCommentService(
	commentRepo port.CommentRepository,
	postRepo port.PostRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *CommentService) CreateComment(
	ctx context.Context,
	postID string,
	parentID *string,
	authorID string,
	content string,
) (*domain.Comment, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if !post.CommentsEnabled {
		return nil, domain.ErrCommentsDisabled
	}
	if len(content) == 0 {
		return nil, domain.ErrEmptyComment
	}
	if len(content) > domain.MaxCommentLength {
		return nil, domain.ErrCommentTooLong
	}
	if parentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, domain.ErrParentNotFound
		}
		depth, err := s.getCommentDepth(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get comment depth: %w", err)
		}
		if depth >= MaxCommentDepth {
			return nil, domain.ErrMaxDepthExceeded
		}
		if parent.PostID != postID {
			return nil, fmt.Errorf("parent comment belongs to different post")
		}
	}

	comment := &domain.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		ParentID:  parentID,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	return comment, nil
}

func (s *CommentService) GetCommentByID(ctx context.Context, id string) (*domain.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return comment, nil
}

func (s *CommentService) ListComments(
	ctx context.Context,
	postID string,
	parentID *string,
	cursor string,
	limit int,
) (*port.CommentListResult, error) {
	exists, err := s.postRepo.Exists(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !exists {
		return nil, domain.ErrPostNotFound
	}

	params := port.CommentListParams{
		PostID:   postID,
		ParentID: parentID,
		Cursor:   cursor,
		Limit:    limit,
	}
	result, err := s.commentRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	return result, nil
}

func (s *CommentService) getCommentDepth(ctx context.Context, commentID string) (int, error) {
	depth := 0
	currentID := commentID
	for currentID != "" {
		comment, err := s.commentRepo.GetByID(ctx, currentID)
		if err != nil {
			return 0, err
		}
		if comment.ParentID == nil {
			break
		}
		currentID = *comment.ParentID
		depth++
		if depth >= MaxCommentDepth {
			return 0, domain.ErrMaxDepthExceeded
		}
	}
	return depth, nil
}
