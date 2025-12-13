package postItems

import (
	"fmt"
	"mime/multipart"
	postItemsDao "uttc-hackathon-backend/dao/postItems"
)

type ItemUsecase struct {
	postItemsDao *postItemsDao.ItemDAO
}

func NewItemUsecase(dao *postItemsDao.ItemDAO) *ItemUsecase {
	return &ItemUsecase{postItemsDao: dao}
}

// CreateItemメソッド（ロジックを外部関数に委譲）
func (h *ItemUsecase) CreateItem(title string, explanation string, priceStr string, file multipart.File, fileHeader *multipart.FileHeader, uid string, ifPurchased bool, category string) (map[string]interface{}, []string, error) {

	// 1. 価格の検証と変換を price.go に委譲
	price := priceToInt(priceStr)

	// 2. ファイルI/O処理を file.go に委譲
	imagePath, err := saveUploadedFile(file, fileHeader)
	if err != nil {
		return nil, nil, fmt.Errorf("file processing error: %w", err)
	}

	// 画像URLを配列として保持（現在は1枚のみ対応）
	var imageURLs []string
	if imagePath != "" {
		imageURLs = append(imageURLs, imagePath)
	}

	// 3. DAOの呼び出し（永続化）
	if err := h.postItemsDao.InsertItem(title, price, explanation, imageURLs, uid, ifPurchased, category); err != nil {
		return nil, nil, fmt.Errorf("database error: %w", err)
	}

	// 4. 成功時のレスポンス
	response := map[string]interface{}{
		"message":    "Item Created successfully",
		"image_urls": imageURLs,
	}

	return response, imageURLs, nil
}
