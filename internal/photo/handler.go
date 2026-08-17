package photo

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"photo-memo/internal/memo"
	"strconv"
	"strings"
)

const maxUploadSize = 10 << 20 // 10MB

// 拡張子は小文字化した上でこのマップと突き合わせる
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 拡張子詐称を防ぐため、実データの中身も検査する
var allowedContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func CreateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		memoID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invaild memo id", http.StatusBadRequest)
			return
		}

		// 紐づけ先のメモが存在するか確認する
		_, err = memo.FindbyOne(db, memoID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "memo not found", http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}

		// リクエスト全体のサイズに上限を設ける
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		//以下ファイルの処理を書き進めていく
		file, header, err := r.FormFile("photo")
		if err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "invaild file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 拡張子を検証する
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowedExtensions[ext] {
			http.Error(w, "unsupported file type", http.StatusBadRequest)
			return
		}

		// 拡張子だけでなく、実際の中身(先頭バイト)も検証する
		sniffBuf := make([]byte, 512)
		n, err := file.Read(sniffBuf)
		if err != nil && err != io.EOF {
			http.Error(w, "invaild file", http.StatusBadRequest)
			return
		}
		contentType := http.DetectContentType(sniffBuf[:n])
		if !allowedContentTypes[contentType] {
			http.Error(w, "unsupported file type", http.StatusBadRequest)
			return
		}
		// 中身を読んだ分、先頭に読み戻す
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, "invaild file", http.StatusInternalServerError)
			return
		}

		// ファイル名を抜き出して、ランダムな名前を割り当てる
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

func FindByMemoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		memoID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invaild memo id", http.StatusBadRequest)
			return
		}

		_, err = memo.FindbyOne(db, memoID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "memo not found", http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}

		photos, err := FindByMemoID(db, memoID)
		if err != nil {
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&photos)
	}
}

func ServeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invaild id", http.StatusBadRequest)
			return
		}

		p, err := FindbyOne(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "photo not found", http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}
		http.ServeFile(w, r, p.FilePath)
	}
}

func DeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invaild id", http.StatusBadRequest)
			return
		}

		p, err := FindbyOne(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "photo not found", http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}

		err = Delete(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "photo not found", http.StatusNotFound)
				return
			}
			log.Println(err)
			http.Error(w, "find miss", http.StatusInternalServerError)
			return
		}

		// DBの整合性を優先し、ファイル削除の失敗はログのみに留める
		if err := os.Remove(p.FilePath); err != nil && !os.IsNotExist(err) {
			log.Println(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
