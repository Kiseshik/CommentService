package service

import (
	"context"
	"fmt"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
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

func (s *CommentService) CreateComment(ctx context.Context, params *port.CreateCommentParams) (*domain.Comment, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, err)
	}
	post, err := s.postRepo.GetByID(ctx, params.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if !post.CommentsEnabled {
		return nil, domain.ErrCommentsDisabled
	}
	if params.ParentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *params.ParentID)
		if err != nil {
			return nil, domain.ErrParentNotFound
		}
		if err := s.validateCommentDepth(ctx, *params.ParentID); err != nil {
			return nil, err
		}
		if parent.PostID != params.PostID {
			return nil, domain.ErrParentFromDifferentPost
		}
	}
	comment, err := s.commentRepo.Create(ctx, params)
	if err != nil {
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

func (s *CommentService) ListComments(ctx context.Context, params *port.ListCommentParams) (*port.CommentListResult, error) {
	exists, err := s.postRepo.Exists(ctx, params.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !exists {
		return nil, domain.ErrPostNotFound
	}
	repoParams := port.CommentListParams{
		PostID:   params.PostID,
		ParentID: params.ParentID,
		Cursor:   params.Cursor,
		Limit:    params.Limit,
	}
	result, err := s.commentRepo.List(ctx, repoParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	return result, nil
}

func (s *CommentService) validateCommentDepth(ctx context.Context, commentID string) error {
	depth := 0
	currentID := commentID
	for currentID != "" {
		comment, err := s.commentRepo.GetByID(ctx, currentID)
		if err != nil {
			return err
		}
		if comment.ParentID == nil {
			break
		}
		currentID = *comment.ParentID
		depth++
		if depth >= MaxCommentDepth {
			return domain.ErrMaxDepthExceeded
		}
	}
	return nil
}
