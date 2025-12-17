package postItems

import (
	"database/sql"
)

type ItemDAO struct {
	db *sql.DB
}

func NewItemDAO(db *sql.DB) *ItemDAO {
	return &ItemDAO{db: db}
}

func (d *ItemDAO) InsertItem(title string, price int, explanation string, imageURLs []string, uid string, status string, category string) error {
	// トランザクション開始
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// itemsテーブルに挿入
	query := "INSERT INTO items (title, price, explanation, uid, status, category) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(query, title, price, explanation, uid, status, category)
	if err != nil {
		return err
	}

	// 挿入したアイテムのIDを取得
	itemID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// item_imagesテーブルに画像URLを挿入
	if len(imageURLs) > 0 {
		imageQuery := "INSERT INTO item_images (item_id, image_url) VALUES (?, ?)"
		for _, url := range imageURLs {
			if url != "" {
				_, err := tx.Exec(imageQuery, itemID, url)
				if err != nil {
					return err
				}
			}
		}
	}

	// コミット
	return tx.Commit()
}
