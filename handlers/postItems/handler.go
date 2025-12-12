package postItems

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	r.ParseForm()

	title := r.PostForm.Get("title")
	explanation := r.PostForm.Get("explanation")
	priceStr := r.PostForm.Get("price")
	file, fileHeader, err := r.FormFile("image")
	uid := r.PostForm.Get("sellerUid")
	ifPurchasedStr := r.PostForm.Get("ifpurchased")
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

	var ifPurchased bool
	if ifPurchasedStr == "" {
		ifPurchased = false
	} else {
		var err error
		ifPurchased, err = strconv.ParseBool(ifPurchasedStr)
		if err != nil {
			fmt.Println("Error parsing ifPurchased:", err)
			writeJSONError(w, "Invalid ifpurchased value", http.StatusBadRequest)
			return
		}
	}

	response, imagePath, err := h.postItemsUc.CreateItem(title, priceStr, explanation, file, fileHeader, uid, ifPurchased, category)
	if err != nil {
		fmt.Println("Error creating item:", err)
		writeJSONError(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	fmt.Printf("出品データ保存完了: %s\n", title)

	w.Header().Set("Content-Type", "application/json")
	response = map[string]string{
		"message":   "Item Created successfully",
		"image_url": imagePath,
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
