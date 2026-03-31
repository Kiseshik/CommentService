package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{db: db}
}

const createPostQuery = `
	insert into posts
	(
		id,
		title,
		content,
		author_id,
		comments_enabled,
		created_at,
		updated_at
	)
	values
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)
`

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	_, err := r.db.ExecContext(ctx, createPostQuery,
		post.ID,
		post.Title,
		post.Content,
		post.AuthorID,
		post.CommentsEnabled,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

const getPostByIDQuery = `
	select
		id,
		title,
		content,
		author_id,
		comments_enabled,
		created_at,
		updated_at
	from posts
	where id = $1
`

func (r *PostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post
	err := r.db.GetContext(ctx, &post, getPostByIDQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

const updatePostQuery = `
	update posts
	set
		title = $1,
		content = $2,
		comments_enabled = $3,
		updated_at = $4
	where id = $5
`

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	result, err := r.db.ExecContext(ctx, updatePostQuery,
		post.Title,
		post.Content,
		post.CommentsEnabled,
		post.UpdatedAt,
		post.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrPostNotFound
	}
	return nil
}

const existsPostQuery = `
	select exists(
		select 1
		from posts
		where id = $1
	)
`

func (r *PostRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, existsPostQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check post existence: %w", err)
	}
	return exists, nil
}

const listPostsQuery = `
	select
		id,
		title,
		content,
		author_id,
		comments_enabled,
		created_at,
		updated_at
	from posts
	order by created_at desc
`

func (r *PostRepository) List(ctx context.Context) ([]*domain.Post, error) {
	posts := make([]*domain.Post, 0)
	err := r.db.SelectContext(ctx, &posts, listPostsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	return posts, nil
}
