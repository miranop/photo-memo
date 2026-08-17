package photo

import (
	"database/sql"
	"time"
)

func Create(db *sql.DB, memoID int, filePath string) (Photo, error) {
	now := time.Now()
	result, err := db.Exec(`INSERT INTO photo (memo_id,file_path,created_at) VALUES (?,?,?)`, memoID, filePath, now)
	if err != nil {
		return Photo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Photo{}, err
	}
	photo := Photo{
		ID:        int(id),
		MemoID:    memoID,
		FilePath:  filePath,
		CreatedAt: now,
	}
	return photo, nil
}

func FindByMemoID(db *sql.DB, memoID int) ([]Photo, error) {
	var photos []Photo
	rows, err := db.Query("SELECT id, memo_id, file_path, created_at FROM photo WHERE memo_id = ?", memoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Photo
		err := rows.Scan(&p.ID, &p.MemoID, &p.FilePath, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return photos, nil
}

func FindbyOne(db *sql.DB, id int) (Photo, error) {
	var p Photo
	err := db.QueryRow("SELECT id, memo_id, file_path, created_at FROM photo WHERE id = ?", id).Scan(&p.ID, &p.MemoID, &p.FilePath, &p.CreatedAt)
	if err != nil {
		return Photo{}, err
	}
	return p, nil
}

func Delete(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM photo WHERE id = ?", id)
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
