package memo

import "time"

type Memo struct {
	ID        int
	UserID    int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type PostMemo struct {
	Title string
	Body  string
}
type PatchMemo struct {
	Title *string
	Body  *string
}
