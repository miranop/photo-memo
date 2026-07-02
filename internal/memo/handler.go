package memo

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p PostMemo
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		memo, err := Create(db, p)
		if err != nil {
			http.Error(w, "SQL miss", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&memo)
	}
}
