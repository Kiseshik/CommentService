package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

const createCommentQuery = `
	insert into comments
	(
		id,
		post_id,
		parent_id,
		content,
		author_id,
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

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	_, err := r.db.ExecContext(ctx, createCommentQuery,
		comment.ID,
		comment.PostID,
		comment.ParentID,
		comment.Content,
		comment.AuthorID,
		comment.CreatedAt,
		comment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

const getCommentByIDQuery = `
	select
		id,
		post_id,
		parent_id,
		content,
		author_id,
		created_at,
		updated_at
	from comments
	where id = $1
`

func (r *CommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.GetContext(ctx, &comment, getCommentByIDQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCommentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	return &comment, nil
}

const updateCommentQuery = `
	update comments
	set
		content = $1,
		updated_at = $2
	where id = $3
`

func (r *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	result, err := r.db.ExecContext(ctx, updateCommentQuery,
		comment.Content,
		comment.UpdatedAt,
		comment.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

const deleteCommentQuery = `
	delete from comments
	where id = $1
`

func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, deleteCommentQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

const existsCommentQuery = `
	select exists(
		select 1
		from comments
		where id = $1
	)
`

func (r *CommentRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, existsCommentQuery, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check comment existence: %w", err)
	}
	return exists, nil
}

const countCommentsByPostQuery = `
	select count(*)
	from comments
	where post_id = $1
`

func (r *CommentRepository) CountByPost(ctx context.Context, postID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, countCommentsByPostQuery, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}
	return count, nil
}

const listCommentsQuery = `
	select
		id,
		post_id,
		parent_id,
		content,
		author_id,
		created_at,
		updated_at
	from comments
	where post_id = $1
`

func (r *CommentRepository) List(ctx context.Context, params port.CommentListParams) (*port.CommentListResult, error) {
	var nextCursor string

	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}

	query := listCommentsQuery
	args := []interface{}{params.PostID}
	argIdx := 2

	//собираем динамически
	if params.ParentID != nil {
		query += fmt.Sprintf(" and parent_id = $%d", argIdx)
		args = append(args, *params.ParentID)
		argIdx++
	} else {
		query += " and parent_id is null"
	}
	query += " order by created_at, id"
	if params.Cursor != "" {
		query += fmt.Sprintf(" and (created_at, id) > (select created_at, id from comments where id = $%d)", argIdx)
		args = append(args, params.Cursor)
		argIdx++
	}
	query += fmt.Sprintf(" limit $%d", argIdx)
	args = append(args, params.Limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	defer rows.Close()

	comments := make([]*domain.Comment, 0, params.Limit+1)
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Content,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, &comment)
	}

	hasNextPage := len(comments) > params.Limit
	if hasNextPage {
		comments = comments[:params.Limit]
	}
	if hasNextPage && len(comments) > 0 {
		nextCursor = comments[len(comments)-1].ID
	}
	return &port.CommentListResult{
		Comments:    comments,
		NextCursor:  nextCursor,
		HasNextPage: hasNextPage,
	}, nil
}
