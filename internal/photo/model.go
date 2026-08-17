package photo

import "time"

type Photo struct {
	ID        int
	MemoID    int
	FilePath  string
	CreatedAt time.Time
}
