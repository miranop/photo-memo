package memo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
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

func FindAllHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memos, err := FindAll(db)
		if err != nil {
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&memos)
	}
}
func FindbyOneHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		m, err := FindbyOne(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "memo not find", http.StatusNotFound)
				return
			}
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&m)
	}
}

func PatchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		var p PatchMemo
		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		m, err := Patch(db, id, p)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "memo update failed", http.StatusNotFound)
				return
			}
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&m)
	}
}
func DeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		err = Delete(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
