package memory

import (
	"context"
	"testing"
	"time"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentRepository_CreateAndGet(t *testing.T) {
	repo := NewCommentRepository()
	tests := []struct {
		name    string
		comment *domain.Comment
		wantErr bool
	}{
		{
			name: "success create root comment",
			comment: &domain.Comment{
				ID:        "comment-1",
				PostID:    "post-1",
				ParentID:  nil,
				Content:   "test comment",
				AuthorID:  "user-1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "success create reply comment",
			comment: &domain.Comment{
				ID:        "comment-2",
				PostID:    "post-1",
				ParentID:  stringPtr("comment-1"),
				Content:   "reply comment",
				AuthorID:  "user-2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(context.Background(), tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			got, err := repo.GetByID(context.Background(), tt.comment.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.comment.ID, got.ID)
			assert.Equal(t, tt.comment.Content, got.Content)
			assert.Equal(t, tt.comment.PostID, got.PostID)
			if tt.comment.ParentID != nil {
				assert.Equal(t, *tt.comment.ParentID, *got.ParentID)
			} else {
				assert.Nil(t, got.ParentID)
			}
		})
	}
}

func TestCommentRepository_ListRootComments(t *testing.T) {
	repo := NewCommentRepository()
	parentID := "c1"
	comments := []*domain.Comment{
		{ID: "c1", PostID: "post-1", ParentID: nil, Content: "root 1", AuthorID: "user-1", CreatedAt: time.Now()},
		{ID: "c2", PostID: "post-1", ParentID: nil, Content: "root 2", AuthorID: "user-1", CreatedAt: time.Now().Add(time.Second)},
		{ID: "c3", PostID: "post-1", ParentID: &parentID, Content: "reply", AuthorID: "user-2", CreatedAt: time.Now()},
	}
	for _, c := range comments {
		_ = repo.Create(context.Background(), c)
	}

	tests := []struct {
		name           string
		params         port.CommentListParams
		wantCount      int
		wantFirstID    string
		wantHasNext    bool
		wantNextCursor string
	}{
		{
			name: "get root comments with limit 10",
			params: port.CommentListParams{
				PostID:   "post-1",
				ParentID: nil,
				Cursor:   "",
				Limit:    10,
			},
			wantCount:      2,
			wantFirstID:    "c1",
			wantHasNext:    false,
			wantNextCursor: "",
		},
		{
			name: "get root comments with limit 1",
			params: port.CommentListParams{
				PostID:   "post-1",
				ParentID: nil,
				Cursor:   "",
				Limit:    1,
			},
			wantCount:      1,
			wantFirstID:    "c1",
			wantHasNext:    true,
			wantNextCursor: "c1",
		},
		{
			name: "get root comments with cursor",
			params: port.CommentListParams{
				PostID:   "post-1",
				ParentID: nil,
				Cursor:   "c1",
				Limit:    10,
			},
			wantCount:      1,
			wantFirstID:    "c2",
			wantHasNext:    false,
			wantNextCursor: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.List(context.Background(), tt.params)
			require.NoError(t, err)
			assert.Len(t, result.Comments, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, tt.wantFirstID, result.Comments[0].ID)
			}
			assert.Equal(t, tt.wantHasNext, result.HasNextPage)
			assert.Equal(t, tt.wantNextCursor, result.NextCursor)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
