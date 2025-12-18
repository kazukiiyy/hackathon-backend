package purchaseItem

import (
	"database/sql"
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
		return err
	}
	defer tx.Rollback()

	// itemsテーブルのstatusとbuyer_addressを更新
	// statusが'listed'または'purchased'の場合に更新（重複イベントに対応）
	updateQuery := "UPDATE items SET status = 'purchased', buyer_address = ? WHERE id = ? AND status IN ('listed', 'purchased')"
	result, err := tx.Exec(updateQuery, buyerAddress, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// 既に'completed'または'cancelled'の場合は更新しない（正常な状態）
		// ただし、itemIDが存在しない場合はエラーを返す
		var currentStatus string
		checkQuery := "SELECT status FROM items WHERE id = ?"
		err := tx.QueryRow(checkQuery, itemID).Scan(&currentStatus)
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		if err != nil {
			return err
		}
		// 既に'completed'または'cancelled'の場合は正常終了
		if currentStatus == "completed" || currentStatus == "cancelled" {
			return tx.Commit()
		}
		// その他の場合はエラー
		return sql.ErrNoRows
	}

	// purchasesテーブルに購入情報を挿入（buyer_addressも含む）
	// 重複チェック: 既に同じitem_idとbuyer_addressの組み合わせが存在する場合はスキップ
	insertQuery := `
		INSERT INTO purchases (item_id, buyer_uid, buyer_address) 
		SELECT ?, ?, ? 
		WHERE NOT EXISTS (
			SELECT 1 FROM purchases 
			WHERE item_id = ? AND buyer_address = ?
		)
	`
	_, err = tx.Exec(insertQuery, itemID, buyerUID, buyerAddress, itemID, buyerAddress)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	query := `
		SELECT i.id, i.title, i.price, i.explanation, i.uid, i.category, p.purchased_at
		FROM purchases p
		JOIN items i ON p.item_id = i.id
		WHERE p.buyer_uid = ?
		ORDER BY p.purchased_at DESC
	`
	rows, err := d.db.Query(query, buyerUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*PurchasedItem
	for rows.Next() {
		var item PurchasedItem
		err := rows.Scan(&item.ID, &item.Title, &item.Price, &item.Explanation, &item.UID, &item.Category, &item.PurchasedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 各アイテムの画像URLを取得
	for _, item := range items {
		urls, err := d.getImageURLsForItem(item.ID)
		if err != nil {
			return nil, err
		}
		item.ImageURLs = urls
	}

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
