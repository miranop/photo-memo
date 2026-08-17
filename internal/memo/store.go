package memo

import (
	"database/sql"
	"time"
)

func Create(db *sql.DB, p PostMemo) (Memo, error) {
	userID := 1 //TODO　あとで認証をちゃんとやる
	now := time.Now()
	result, err := db.Exec(
		`INSERT INTO memo (user_id, title, body,created_at, updated_at) VALUES(?, ?, ?, ?, ?)`,
		userID, p.Title, p.Body, now, now,
	)
	if err != nil {
		return Memo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Memo{}, err
	}
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

func FindAll(db *sql.DB) ([]Memo, error) {
	var memos []Memo
	//全体を読み出すためのSQLの用意
	rows, err := db.Query("SELECT id,user_id,title,body,created_at,updated_at FROM memo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Memo
		err := rows.Scan(&m.ID, &m.UserID, &m.Title, &m.Body, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		memos = append(memos, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return memos, nil
}

func FindbyOne(db *sql.DB, id int) (Memo, error) {
	var m Memo
	err := db.QueryRow("SELECT id, user_id, title, body, created_at, updated_at FROM memo WHERE id = ?", id).Scan(&m.ID, &m.UserID, &m.Title, &m.Body, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return Memo{}, err
	}
	return m, nil
}

func Patch(db *sql.DB, id int, p PatchMemo) (Memo, error) {
	now := time.Now()

	if p.Title != nil {
		_, err := db.Exec(
			"UPDATE memo SET title = ?, updated_at = ? WHERE id = ?",
			*p.Title, now, id,
		)
		if err != nil {
			return Memo{}, err
		}
	}
	if p.Body != nil {
		_, err := db.Exec(
			"UPDATE memo SET body = ?, updated_at = ? WHERE id = ?",
			*p.Body, now, id,
		)
		if err != nil {
			return Memo{}, err
		}
	}

	return FindbyOne(db, id)
}

func Delete(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM memo WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // 該当データがなかった、という扱いにする
	}
	return nil
}
