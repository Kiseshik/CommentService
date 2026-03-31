package utils

import (
	"errors"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
)

func PaginateComments(comments []*domain.Comment, params port.CommentListParams) (*port.CommentListResult, error) {
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
			return nil, errors.New("invalid cursor: comment not found")
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

func PaginatePosts(posts []*domain.Post, params port.PostListParams) (*port.PostListResult, error) {
	start := 0
	if params.Cursor != "" {
		found := false
		for i, p := range posts {
			if p.ID == params.Cursor {
				start = i + 1
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("invalid cursor: post not found")
		}
	}
	end := start + params.Limit
	if end > len(posts) {
		end = len(posts)
	}
	res := &port.PostListResult{
		Posts:       posts[start:end],
		HasNextPage: end < len(posts),
	}
	if res.HasNextPage && len(res.Posts) > 0 {
		res.NextCursor = res.Posts[len(res.Posts)-1].ID
	}
	return res, nil
}

// TODO: В будущем можно будет объединить PaginateComments и PaginatePosts в одну универсальную функцию с использованием дженериков
