package domain

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidCursor = errors.New("invalid cursor")
)

//TODO: обновить во всем проекте ошибки
