package likes

import (
	"database/sql"
	"time"
)

type Like struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	UID       string    `json:"uid"`
	CreatedAt time.Time `json:"created_at"`
}

// LikeDAOInterface はモック化のためのインターフェース
type LikeDAOInterface interface {
	AddLike(itemID int, uid string) error
	RemoveLike(itemID int, uid string) error
	IsLiked(itemID int, uid string) (bool, error)
	GetLikeCount(itemID int) (int, error)
	GetLikedItemsByUser(uid string) ([]int, error)
}

type LikeDAO struct {
	db *sql.DB
}

func NewLikeDAO(db *sql.DB) *LikeDAO {
	return &LikeDAO{db: db}
}

// いいねを追加（トランザクションでlike_countも更新）
func (d *LikeDAO) AddLike(itemID int, uid string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// likesテーブルに追加
	_, err = tx.Exec("INSERT INTO likes (item_id, uid) VALUES (?, ?)", itemID, uid)
	if err != nil {
		return err
	}

	// itemsのlike_countを+1
	_, err = tx.Exec("UPDATE items SET like_count = like_count + 1 WHERE id = ?", itemID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// いいねを削除（トランザクションでlike_countも更新）
func (d *LikeDAO) RemoveLike(itemID int, uid string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// likesテーブルから削除
	result, err := tx.Exec("DELETE FROM likes WHERE item_id = ? AND uid = ?", itemID, uid)
	if err != nil {
		return err
	}

	// 削除された場合のみlike_countを-1
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		_, err = tx.Exec("UPDATE items SET like_count = GREATEST(like_count - 1, 0) WHERE id = ?", itemID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ユーザーが特定の商品にいいねしているか確認
func (d *LikeDAO) IsLiked(itemID int, uid string) (bool, error) {
	query := "SELECT COUNT(*) FROM likes WHERE item_id = ? AND uid = ?"
	var count int
	err := d.db.QueryRow(query, itemID, uid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 商品のいいね数を取得
func (d *LikeDAO) GetLikeCount(itemID int) (int, error) {
	query := "SELECT COUNT(*) FROM likes WHERE item_id = ?"
	var count int
	err := d.db.QueryRow(query, itemID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ユーザーがいいねした商品一覧を取得
func (d *LikeDAO) GetLikedItemsByUser(uid string) ([]int, error) {
	query := "SELECT item_id FROM likes WHERE uid = ? ORDER BY created_at DESC"
	rows, err := d.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemIDs []int
	for rows.Next() {
		var itemID int
		if err := rows.Scan(&itemID); err != nil {
			return nil, err
		}
		itemIDs = append(itemIDs, itemID)
	}
	return itemIDs, rows.Err()
}
