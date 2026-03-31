package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrInvalidInput = errors.New("invalid input")

	ErrPostNotFound      = errors.New("post not found")
	ErrPostAlreadyExists = errors.New("post already exists")

	ErrCommentNotFound      = errors.New("comment not found")
	ErrEmptyComment         = errors.New("comment cannot be empty")
	ErrCommentAlreadyExists = errors.New("comment already exists")
	ErrCommentTooLong       = errors.New("comment exceeds maximum length of 2000 characters")
	ErrMaxDepthExceeded     = errors.New("max depth exceeded")
	ErrCommentsDisabled     = errors.New("comments are disabled for this post")
	ErrParentNotFound       = errors.New("parent comment not found")
)

//TODO: обновить во всем проекте ошибки
