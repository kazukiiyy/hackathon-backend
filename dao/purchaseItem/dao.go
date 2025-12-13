package purchaseItem

import (
	"database/sql"
)

type PurchaseDAO struct {
	db *sql.DB
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
