package photo

import (
	"database/sql"
	"time"
)

func Create(db *sql.DB, memoID int, filePath string) (Photo, error) {
	now := time.Now()
	result, err := db.Exec(`INSERT INTO photo (memo_id,file_path) VALUES (?,?)`, memoID, filePath)
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
