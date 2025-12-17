package postItems

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uttc-hackathon-backend/usecase/postItems"
)

type ItemHandler struct {
	postItemsUc *postItems.ItemUsecase
}

func NewItemHandler(u *postItems.ItemUsecase) *ItemHandler {
	return &ItemHandler{postItemsUc: u}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20) // 10MB max

	title := r.PostForm.Get("title")
	explanation := r.PostForm.Get("explanation")
	priceStr := r.PostForm.Get("price")
	file, fileHeader, err := r.FormFile("image")
	uid := r.PostForm.Get("sellerUid")
	status := r.PostForm.Get("status")
	category := r.PostForm.Get("category")

	if err != nil && err != http.ErrMissingFile {
		// ファイル取得自体の内部エラー
		writeJSONError(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// statusが空の場合はデフォルトで"listed"を設定
	if status == "" {
		status = "listed"
	}

	// バリデーション
	if title == "" {
		writeJSONError(w, "title is required", http.StatusBadRequest)
		return
	}
	if uid == "" {
		writeJSONError(w, "uid is required", http.StatusBadRequest)
		return
	}
	if category == "" {
		writeJSONError(w, "category is required", http.StatusBadRequest)
		return
	}
	if priceStr == "" {
		writeJSONError(w, "price is required", http.StatusBadRequest)
		return
	}

	response, imageURLs, err := h.postItemsUc.CreateItem(title, explanation, priceStr, file, fileHeader, uid, status, category)
	if err != nil {
		fmt.Printf("Error creating item - Title: %s, UID: %s, Error: %v\n", title, uid, err)
		// エラーメッセージを詳細に返す（開発環境用）
		errorMsg := fmt.Sprintf("Failed to create item: %v", err)
		writeJSONError(w, errorMsg, http.StatusInternalServerError)
		return
	}

	fmt.Printf("出品データ保存完了: %s\n", title)

	w.Header().Set("Content-Type", "application/json")
	response = map[string]interface{}{
		"message":    "Item Created successfully",
		"image_urls": imageURLs,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeJSONError(w, "JSON encode error", http.StatusInternalServerError)
	}

}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
