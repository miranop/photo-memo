package memo

import (
	"database/sql"
	"time"
)

func Create(db *sql.DB, p PostMemo) (Memo, error) {
	userID := 1 //TODO　あとで認証をちゃんとやる
	result, err := db.Exec(
		`INSERT INTO memo (user_id, title, body) VALUES(?, ?, ?)`,
		userID, p.Title, p.Body,
	)
	if err != nil {
		return Memo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Memo{}, err
	}
	now := time.Now()
	memo := Memo{
		ID:        int(id),
		UserID:    userID,
		Title:     p.Title,
		Body:      p.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return memo, nil
}
