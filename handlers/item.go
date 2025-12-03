package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"uttc-hackathon-backend/dao"
)

type ItemHandler struct {
	dao *dao.UserDAO
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	title := r.PostForm.Get("title")
	explanation := r.PostForm.Get("explanation")
	priceStr := r.PostForm.Get("price")

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	var imagePath string

	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		uploadDir := "./uploads"
		if err := os.MkdirALL(uploadDir, os.MadePerm); err != nil {
			http.Error(w, "Could not create directory", http.StatusInternalServerError)
			return
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		filepath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "Could not savd file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}
		imagePath = filepath
		fmt.Printf("画像が保存されました : %s\n", imagePath)
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}

	err = h.dao.InsertItem(title, price, explanation, imagePath)
	if err != nil {
		fmt.Printf("DB Error: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return

	}
	fmt.Printf("出品データ保存完了: %s\n", title)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message":   "Item Created successfully",
		"image_url": imagePath,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}

}
