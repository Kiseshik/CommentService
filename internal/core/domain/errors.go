package domain

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidCursor = errors.New("invalid cursor")
)

//todo затестить как возвращаются ошибки с постгри, скорее всего есть дубликаты, посмотреть по всему проекту ниже мемори репы
//repository/postgres/
//app/
//controllers/api/api.go
