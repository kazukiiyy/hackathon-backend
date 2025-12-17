package postItems

import (
	"database/sql"
	"fmt"
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

	// priceを文字列に変換（データベースではVARCHARとして保存）
	priceStr := fmt.Sprintf("%d", price)

	// itemsテーブルに挿入
	query := "INSERT INTO items (title, price, explanation, uid, status, category) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(query, title, priceStr, explanation, uid, status, category)
	if err != nil {
		return fmt.Errorf("failed to insert item into database: %w", err)
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
					return fmt.Errorf("failed to insert image URL into database: %w", err)
				}
			}
		}
	}

	// コミット
	return tx.Commit()
}

// UpdateChainItemID は既存の商品にchain_item_idを関連付ける
func (d *ItemDAO) UpdateChainItemID(itemID int64, chainItemID int64, sellerAddress string, tokenID int64) error {
	query := "UPDATE items SET chain_item_id = ?, seller_address = ?, token_id = ? WHERE id = ?"
	_, err := d.db.Exec(query, chainItemID, sellerAddress, tokenID, itemID)
	return err
}

// FindItemByUidAndTitle はuidとtitleで商品を検索（chain_item_idを関連付けるため）
func (d *ItemDAO) FindItemByUidAndTitle(uid string, title string) (int64, error) {
	query := "SELECT id FROM items WHERE uid = ? AND title = ? AND chain_item_id IS NULL ORDER BY created_at DESC LIMIT 1"
	var itemID int64
	err := d.db.QueryRow(query, uid, title).Scan(&itemID)
	if err != nil {
		return 0, err
	}
	return itemID, nil
}

// FindItemByChainItemID はchain_item_idで商品を検索
func (d *ItemDAO) FindItemByChainItemID(chainItemID int64) (int64, error) {
	query := "SELECT id FROM items WHERE chain_item_id = ? LIMIT 1"
	var itemID int64
	err := d.db.QueryRow(query, chainItemID).Scan(&itemID)
	if err != nil {
		return 0, err
	}
	return itemID, nil
}

// InsertItemWithChainID はchain_item_idを含めて商品を挿入
func (d *ItemDAO) InsertItemWithChainID(title string, price int, explanation string, imageURLs []string, uid string, status string, category string, chainItemID int64, sellerAddress string, tokenID int64) error {
	// トランザクション開始
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// priceを文字列に変換（データベースではVARCHARとして保存）
	priceStr := fmt.Sprintf("%d", price)

	// itemsテーブルに挿入（chain_item_idを含む）
	query := "INSERT INTO items (title, price, explanation, uid, status, category, chain_item_id, seller_address, token_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(query, title, priceStr, explanation, uid, status, category, chainItemID, sellerAddress, tokenID)
	if err != nil {
		// より詳細なエラーメッセージを返す
		return fmt.Errorf("failed to insert item into items table (title=%s, chain_item_id=%d, uid=%s): %w", title, chainItemID, uid, err)
	}

	// 挿入したアイテムのIDを取得
	itemID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// item_imagesテーブルに画像URLを挿入
	if len(imageURLs) > 0 {
		// chain_item_idカラムがある場合は含める、ない場合はitem_idのみ
		imageQuery := "INSERT INTO item_images (item_id, image_url, chain_item_id) VALUES (?, ?, ?)"
		for _, url := range imageURLs {
			if url != "" {
				_, err := tx.Exec(imageQuery, itemID, url, chainItemID)
				if err != nil {
					// chain_item_idカラムがない可能性があるので、item_idのみで再試行
					imageQueryFallback := "INSERT INTO item_images (item_id, image_url) VALUES (?, ?)"
					_, err2 := tx.Exec(imageQueryFallback, itemID, url)
					if err2 != nil {
						return fmt.Errorf("failed to insert image URL: %w (original error: %v)", err2, err)
					}
				}
			}
		}
	}

	// コミット
	return tx.Commit()
}
