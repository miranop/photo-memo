package main

import (
	"net/http"
	"photo-memo/internal/db"
	"photo-memo/internal/memo"
)

func main() {
	database := db.DbInit("./photo-memo.db")
	defer database.Close()

	mux := http.NewServeMux()
	mux.Handle("POST /memos", memo.CreateHandler(database))
	http.ListenAndServe(":8080", mux)
}
