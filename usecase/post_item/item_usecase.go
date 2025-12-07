package post_item

import (
	"fmt"
	"mime/multipart"
	"uttc-hackathon-backend/dao"
)

type ItemUsecase struct {
	dao *dao.ItemDAO
}

func NewItemUsecase(dao *dao.ItemDAO) *ItemUsecase {
	return &ItemUsecase{dao: dao}
}

// CreateItemメソッド（ロジックを外部関数に委譲）
func (uc *ItemUsecase) CreateItem(title string, explanation string, priceStr string, file multipart.File, fileHeader *multipart.FileHeader) (map[string]string, string, error) {

	// 1. 価格の検証と変換を price.go に委譲
	price := priceToInt(priceStr)

	// 2. ファイルI/O処理を file.go に委譲
	imagePath, err := saveUploadedFile(file, fileHeader)
	if err != nil {
		return nil, "", fmt.Errorf("file processing error: %w", err)
	}

	// 3. DAOの呼び出し（永続化）
	if err := uc.dao.InsertItem(title, price, explanation, imagePath); err != nil {
		return nil, "", fmt.Errorf("database error: %w", err)
	}

	// 4. 成功時のレスポンス
	response := map[string]string{
		"message":   "Item Created successfully",
		"image_url": imagePath,
	}

	return response, imagePath, nil
}
