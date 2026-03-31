package service

import (
	"context"
	"testing"

	"github.com/Kiseshik/CommentService.git/internal/adapters/repository/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostService_CreateAndGet(t *testing.T) {
	postRepo := memory.NewPostRepository()
	postService := NewPostService(postRepo)
	ctx := context.Background()

	post, err := postService.CreatePost(ctx, "Test Title", "Test Content", "user-1", true)
	require.NoError(t, err)
	require.NotNil(t, post)
	assert.NotEmpty(t, post.ID)
	assert.Equal(t, "Test Title", post.Title)
	assert.Equal(t, "Test Content", post.Content)
	assert.Equal(t, "user-1", post.AuthorID)
	assert.True(t, post.CommentsEnabled)
	assert.NotZero(t, post.CreatedAt)
	assert.NotZero(t, post.UpdatedAt)

	got, err := postService.GetPostByID(ctx, post.ID)
	require.NoError(t, err)
	assert.Equal(t, post.ID, got.ID)
	assert.Equal(t, post.Title, got.Title)

	_, err = postService.GetPostByID(ctx, "not-exists")
	assert.Error(t, err)
}

func TestPostService_ListPosts(t *testing.T) {
	postRepo := memory.NewPostRepository()
	postService := NewPostService(postRepo)
	ctx := context.Background()

	_, err := postService.CreatePost(ctx, "Post 1", "Content 1", "user-1", true)
	require.NoError(t, err)
	_, err = postService.CreatePost(ctx, "Post 2", "Content 2", "user-2", true)
	require.NoError(t, err)
	_, err = postService.CreatePost(ctx, "Post 3", "Content 3", "user-1", false)
	require.NoError(t, err)
	list, err := postService.ListPosts(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 3)

	titles := make(map[string]bool)
	for _, p := range list {
		titles[p.Title] = true
	}
	assert.True(t, titles["Post 1"])
	assert.True(t, titles["Post 2"])
	assert.True(t, titles["Post 3"])
}

func TestPostService_ToggleComments(t *testing.T) {
	postRepo := memory.NewPostRepository()
	postService := NewPostService(postRepo)
	ctx := context.Background()

	post, err := postService.CreatePost(ctx, "Test", "Content", "user-1", true)
	require.NoError(t, err)
	assert.True(t, post.CommentsEnabled)

	err = postService.ToggleComments(ctx, post.ID, false)
	require.NoError(t, err)
	got, err := postService.GetPostByID(ctx, post.ID)
	require.NoError(t, err)
	assert.False(t, got.CommentsEnabled)

	err = postService.ToggleComments(ctx, post.ID, true)
	require.NoError(t, err)
	got, err = postService.GetPostByID(ctx, post.ID)
	require.NoError(t, err)
	assert.True(t, got.CommentsEnabled)

	err = postService.ToggleComments(ctx, "not-exists", true)
	assert.Error(t, err)
}

func TestPostService_Exists(t *testing.T) {
	postRepo := memory.NewPostRepository()
	postService := NewPostService(postRepo)
	ctx := context.Background()

	post, err := postService.CreatePost(ctx, "Test", "Content", "user-1", true)
	require.NoError(t, err)
	exists, err := postService.Exists(ctx, post.ID)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = postService.Exists(ctx, "not-exists")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestPostService_Validation(t *testing.T) {
	postRepo := memory.NewPostRepository()
	postService := NewPostService(postRepo)
	ctx := context.Background()

	tests := []struct {
		name    string
		title   string
		content string
		author  string
		wantErr bool
	}{
		{"valid post", "Title", "Content", "user-1", false},
		{"empty title", "", "Content", "user-1", true},
		{"empty content", "Title", "", "user-1", true},
		{"empty author", "Title", "Content", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := postService.CreatePost(ctx, tt.title, tt.content, tt.author, true)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
