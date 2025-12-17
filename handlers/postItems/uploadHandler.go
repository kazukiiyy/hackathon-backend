package postItems

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// UploadImage は画像のみをアップロードしてURLを返す
func (h *ItemHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10MB max

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		writeJSONError(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// 画像をアップロードしてURLを取得
	imagePath, err := h.postItemsUc.UploadImage(file, fileHeader)
	if err != nil {
		fmt.Printf("Error uploading image: %v\n", err)
		writeJSONError(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// フルURLを構築
	// 環境変数からベースURLを取得、なければリクエストから構築
	baseURL := os.Getenv("BACKEND_BASE_URL")
	if baseURL == "" {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, r.Host)
	}
	imageURL := fmt.Sprintf("%s/%s", baseURL, imagePath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"image_url": imageURL,
		"image_urls": []string{imageURL},
	})
}

