package purchaseItem

import (
	"database/sql"
	"time"
)

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

func (d *PurchaseDAO) UpdatePurchaseStatus(itemID int, buyerUID string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// itemsテーブルのifPurchasedを更新
	updateQuery := "UPDATE items SET ifPurchased = true WHERE id = ? AND ifPurchased = false"
	result, err := tx.Exec(updateQuery, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// purchasesテーブルに購入情報を挿入
	insertQuery := "INSERT INTO purchases (item_id, buyer_uid) VALUES (?, ?)"
	_, err = tx.Exec(insertQuery, itemID, buyerUID)
	if err != nil {
		return err
	}

	return tx.Commit()
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
