package utils

import (
	"sort"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

func SortComments(comments []*domain.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		if comments[i].CreatedAt.Equal(comments[j].CreatedAt) {
			return comments[i].ID < comments[j].ID
		}
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})
}

func SortPosts(posts []*domain.Post) {
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].CreatedAt.Equal(posts[j].CreatedAt) {
			return posts[i].ID < posts[j].ID
		}
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})
}

// TODO: В будущем можно будет объединить SortComments и SortPosts в одну универсальную функцию с использованием дженериков
