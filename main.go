package main

import (
	"net/http"
	"photo-memo/internal/db"
	"photo-memo/internal/memo"
	"photo-memo/internal/photo"
)

func main() {
	database := db.DbInit("./photo-memo.db")
	defer database.Close()

	mux := http.NewServeMux()
	mux.Handle("POST /memos", memo.CreateHandler(database))
	mux.Handle("GET /memos", memo.FindAllHandler(database))
	mux.Handle("GET /memos/{id}", memo.FindbyOneHandler(database))
	mux.Handle("PATCH /memos/{id}", memo.PatchHandler(database))
	mux.Handle("DELETE /memos/{id}", memo.DeleteHandler(database))
	mux.Handle("POST /memos/{id}/photos", photo.CreateHandler(database))
	http.ListenAndServe(":8080", mux)
}
