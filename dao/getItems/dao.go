package getItems

import (
	"database/sql"
	"strconv"
	"time"
)

type Item struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	Explanation string    `json:"explanation"`
	ImageURLs   []string  `json:"image_urls"`
	UID         string    `json:"uid"`
	Status      string    `json:"status"`
	Category    string    `json:"category"`
	LikeCount   int       `json:"like_count"`
	CreatedAt   time.Time `json:"created_at"`
	ChainItemID *int64    `json:"chain_item_id,omitempty"`
	IfPurchased bool      `json:"ifPurchased"`
}

type ItemDAO struct {
	db *sql.DB
}

func NewItemDAO(db *sql.DB) *ItemDAO {
	return &ItemDAO{db: db}
}

func (d *ItemDAO) getImageURLsForItem(itemID int) ([]string, error) {
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

func (d *ItemDAO) GetItemsByCategory(category string, page, limit int) ([]*Item, error) {
	offset := (page - 1) * limit
	query := "SELECT id, title, price, explanation, uid, status, category, like_count, created_at, chain_item_id FROM items WHERE category = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := d.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		var chainItemID sql.NullInt64
		var priceStr string
		err := rows.Scan(&item.ID, &item.Title, &priceStr, &item.Explanation, &item.UID, &item.Status, &item.Category, &item.LikeCount, &item.CreatedAt, &chainItemID)
		if err != nil {
			return nil, err
		}
		// priceを文字列からintに変換
		priceInt, err := strconv.Atoi(priceStr)
		if err != nil {
			return nil, err
		}
		item.Price = priceInt
		if chainItemID.Valid {
			val := chainItemID.Int64
			item.ChainItemID = &val
		}
		item.IfPurchased = item.Status == "purchased" || item.Status == "completed"
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

func (d *ItemDAO) GetItemByID(id int) (*Item, error) {
	query := "SELECT id, title, price, explanation, uid, status, category, like_count, created_at, chain_item_id FROM items WHERE id = ?"
	row := d.db.QueryRow(query, id)

	var item Item
	var chainItemID sql.NullInt64
	var priceStr string
	err := row.Scan(&item.ID, &item.Title, &priceStr, &item.Explanation, &item.UID, &item.Status, &item.Category, &item.LikeCount, &item.CreatedAt, &chainItemID)
	if err != nil {
		return nil, err
	}
	// priceを文字列からintに変換
	priceInt, err := strconv.Atoi(priceStr)
	if err != nil {
		return nil, err
	}
	item.Price = priceInt
	if chainItemID.Valid {
		val := chainItemID.Int64
		item.ChainItemID = &val
	}
	item.IfPurchased = item.Status == "purchased" || item.Status == "completed"

	// 画像URLを取得
	urls, err := d.getImageURLsForItem(item.ID)
	if err != nil {
		return nil, err
	}
	item.ImageURLs = urls

	return &item, nil
}

// 新着商品を取得
func (d *ItemDAO) GetLatestItems(limit int) ([]*Item, error) {
	query := "SELECT id, title, price, explanation, uid, status, category, like_count, created_at, chain_item_id FROM items ORDER BY created_at DESC LIMIT ?"
	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		var chainItemID sql.NullInt64
		var priceStr string
		err := rows.Scan(&item.ID, &item.Title, &priceStr, &item.Explanation, &item.UID, &item.Status, &item.Category, &item.LikeCount, &item.CreatedAt, &chainItemID)
		if err != nil {
			return nil, err
		}
		// priceを文字列からintに変換
		priceInt, err := strconv.Atoi(priceStr)
		if err != nil {
			return nil, err
		}
		item.Price = priceInt
		if chainItemID.Valid {
			val := chainItemID.Int64
			item.ChainItemID = &val
		}
		item.IfPurchased = item.Status == "purchased" || item.Status == "completed"
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, item := range items {
		urls, err := d.getImageURLsForItem(item.ID)
		if err != nil {
			return nil, err
		}
		item.ImageURLs = urls
	}

	return items, nil
}

func (d *ItemDAO) GetItemsByUid(uid string) ([]*Item, error) {
	query := "SELECT id, title, price, explanation, uid, status, category, like_count, created_at, chain_item_id FROM items WHERE uid = ? ORDER BY created_at DESC"
	rows, err := d.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		var chainItemID sql.NullInt64
		var priceStr string
		err := rows.Scan(&item.ID, &item.Title, &priceStr, &item.Explanation, &item.UID, &item.Status, &item.Category, &item.LikeCount, &item.CreatedAt, &chainItemID)
		if err != nil {
			return nil, err
		}
		// priceを文字列からintに変換
		priceInt, err := strconv.Atoi(priceStr)
		if err != nil {
			return nil, err
		}
		item.Price = priceInt
		if chainItemID.Valid {
			val := chainItemID.Int64
			item.ChainItemID = &val
		}
		item.IfPurchased = item.Status == "purchased" || item.Status == "completed"
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
