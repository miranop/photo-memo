package photo

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func CreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		memoID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invaild memo id", http.StatusBadRequest)
			return
		}
		//以下ファイルの処理を書き進めていく
		file, header, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "invaild file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		// ファイル名を抜き出して、ランダムな名前を割り当てる
		ext := filepath.Ext(header.Filename)
		randomBytes := make([]byte, 16)
		_, err = rand.Read(randomBytes)
		if err != nil {
			http.Error(w, "failed to generate filename", http.StatusInternalServerError)
			return
		}
		fileName := hex.EncodeToString(randomBytes) + ext

		//ディレクトリを実際に作成し、画像を保存する
		err = os.MkdirAll("photos", 0755)
		if err != nil {
			http.Error(w, "failed to create directory", http.StatusInternalServerError)
			return
		}
		savePath := filepath.Join("photos", fileName)
		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}
		photo, err := Create(db, memoID, savePath)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to save photo record", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&photo)
	}
}
