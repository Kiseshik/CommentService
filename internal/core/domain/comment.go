package domain

import "time"

type Comment struct {
	ID        string    `db:"id"`
	PostID    string    `db:"post_id"`
	Content   string    `db:"content"`
	AuthorID  string    `db:"author_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
