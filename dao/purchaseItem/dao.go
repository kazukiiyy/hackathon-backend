package purchaseItem

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// PurchaseDAOInterface はモック化のためのインターフェース
type PurchaseDAOInterface interface {
	UpdatePurchaseStatus(itemID int, buyerUID string, buyerAddress string) error
	GetPurchasedItems(buyerUID string) ([]*PurchasedItem, error)
	GetUIDByWalletAddress(walletAddress string) (string, error)
}

type PurchaseDAO struct {
	db *sql.DB
}

type PurchasedItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	Explanation string    `json:"explanation"`
	ImageURLs   []string  `json:"image_urls"`
	UID         string    `json:"uid"`
	Category    string    `json:"category"`
	PurchasedAt time.Time `json:"purchased_at"`
}

func NewPurchaseDAO(db *sql.DB) *PurchaseDAO {
	return &PurchaseDAO{db: db}
}

func (d *PurchaseDAO) UpdatePurchaseStatus(itemID int, buyerUID string, buyerAddress string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// itemsテーブルのstatusとbuyer_addressを更新
	// statusが'listed'または'purchased'の場合に更新（重複イベントに対応）
	updateQuery := "UPDATE items SET status = 'purchased', buyer_address = ? WHERE id = ? AND status IN ('listed', 'purchased')"
	result, err := tx.Exec(updateQuery, buyerAddress, itemID)
	if err != nil {
		if strings.Contains(err.Error(), "Unknown column 'buyer_address'") {
			return fmt.Errorf("buyer_address column does not exist in items table. Please run the migration script: add_buyer_address_column_safe.sql. Original error: %w", err)
		}
		return fmt.Errorf("failed to update purchase status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// 更新されなかった場合、既に完了済みか確認
	if rowsAffected == 0 {
		var currentStatus string
		if err := tx.QueryRow("SELECT status FROM items WHERE id = ?", itemID).Scan(&currentStatus); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("item not found: itemID=%d", itemID)
			}
			return fmt.Errorf("failed to check current status: %w", err)
		}
		// 既に完了済みの場合は正常終了
		if currentStatus == "completed" || currentStatus == "cancelled" {
			return tx.Commit()
		}
		return fmt.Errorf("item status is '%s', cannot update to purchased", currentStatus)
	}

	// purchasesテーブルに購入情報を挿入（重複チェック付き）
	var insertQuery string
	if buyerAddress != "" {
		insertQuery = `
			INSERT INTO purchases (item_id, buyer_uid, buyer_address) 
			SELECT ?, ?, ? 
			WHERE NOT EXISTS (
				SELECT 1 FROM purchases 
				WHERE item_id = ? AND buyer_address = ?
			)
		`
		_, err = tx.Exec(insertQuery, itemID, buyerUID, buyerAddress, itemID, buyerAddress)
	} else {
		insertQuery = `
			INSERT INTO purchases (item_id, buyer_uid, buyer_address) 
			SELECT ?, ?, ? 
			WHERE NOT EXISTS (
				SELECT 1 FROM purchases 
				WHERE item_id = ? AND buyer_uid = ? AND (buyer_address IS NULL OR buyer_address = '')
			)
		`
		_, err = tx.Exec(insertQuery, itemID, buyerUID, buyerAddress, itemID, buyerUID)
	}
	if err != nil {
		return fmt.Errorf("failed to insert purchase record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetUIDByWalletAddress はウォレットアドレスからUIDを取得
func (d *PurchaseDAO) GetUIDByWalletAddress(walletAddress string) (string, error) {
	query := "SELECT uid FROM users WHERE wallet_address = ? LIMIT 1"
	var uid string
	err := d.db.QueryRow(query, walletAddress).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}

// 購入した商品一覧を取得
func (d *PurchaseDAO) GetPurchasedItems(buyerUID string) ([]*PurchasedItem, error) {
	log.Printf("[PurchaseDAO] GetPurchasedItems called with buyerUID=%s", buyerUID)
	
	// まず、purchasesテーブルに該当するレコードが存在するか確認
	var purchaseCount int
	countQuery := "SELECT COUNT(*) FROM purchases WHERE buyer_uid = ?"
	err := d.db.QueryRow(countQuery, buyerUID).Scan(&purchaseCount)
	if err != nil {
		log.Printf("[PurchaseDAO] Count query error: %v", err)
	} else {
		log.Printf("[PurchaseDAO] Total purchases found in purchases table: %d", purchaseCount)
	}
	
	// item_idがNULLでない購入レコードの数を確認
	var validPurchaseCount int
	validCountQuery := "SELECT COUNT(*) FROM purchases WHERE buyer_uid = ? AND item_id IS NOT NULL"
	err = d.db.QueryRow(validCountQuery, buyerUID).Scan(&validPurchaseCount)
	if err != nil {
		log.Printf("[PurchaseDAO] Valid count query error: %v", err)
	} else {
		log.Printf("[PurchaseDAO] Purchases with valid item_id: %d", validPurchaseCount)
	}
	
	query := `
		SELECT i.id, i.title, i.price, i.explanation, i.uid, i.category, p.purchased_at
		FROM purchases p
		JOIN items i ON p.item_id = i.id
		WHERE p.buyer_uid = ?
		ORDER BY p.purchased_at DESC
	`
	log.Printf("[PurchaseDAO] Executing query: %s with buyerUID=%s", query, buyerUID)
	rows, err := d.db.Query(query, buyerUID)
	if err != nil {
		log.Printf("[PurchaseDAO] Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []*PurchasedItem
	rowCount := 0
	for rows.Next() {
		var item PurchasedItem
		var purchasedAt sql.NullTime
		err := rows.Scan(&item.ID, &item.Title, &item.Price, &item.Explanation, &item.UID, &item.Category, &purchasedAt)
		if err != nil {
			log.Printf("[PurchaseDAO] Scan error: %v", err)
			return nil, err
		}
		if purchasedAt.Valid {
			item.PurchasedAt = purchasedAt.Time
		} else {
			item.PurchasedAt = time.Now() // デフォルト値
		}
		items = append(items, &item)
		rowCount++
		log.Printf("[PurchaseDAO] Scanned item: ID=%d, Title=%s", item.ID, item.Title)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[PurchaseDAO] Rows error: %v", err)
		return nil, err
	}

	log.Printf("[PurchaseDAO] Found %d purchased items after JOIN", rowCount)

	// 各アイテムの画像URLを取得
	for _, item := range items {
		urls, err := d.getImageURLsForItem(item.ID)
		if err != nil {
			log.Printf("[PurchaseDAO] Error getting image URLs for item %d: %v", item.ID, err)
			return nil, err
		}
		item.ImageURLs = urls
		log.Printf("[PurchaseDAO] Item %d (%s) has %d images", item.ID, item.Title, len(urls))
	}

	log.Printf("[PurchaseDAO] Returning %d items", len(items))
	return items, nil
}

func (d *PurchaseDAO) getImageURLsForItem(itemID int) ([]string, error) {
	query := "SELECT image_url FROM item_images WHERE item_id = ? ORDER BY id"
	rows, err := d.db.Query(query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, rows.Err()
}
