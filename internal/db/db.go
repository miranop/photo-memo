package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DbInit(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS memo (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    INTEGER,
			title      TEXT NOT NULL,
			body       TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME

		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS photo (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			memo_id     INTEGER NOT NULL,
			file_path   TEXT NOT NULL,
			created_at  DATETIME,
			FOREIGN KEY (memo_id) REFERENCES memo(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
