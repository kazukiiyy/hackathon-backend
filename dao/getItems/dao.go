package getItems

import (
	"database/sql"
	"time"
)

type Item struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	Explanation string    `json:"explanation"`
	ImageURL    string    `json:"image_url"`
	UID         string    `json:"uid"`
	IfPurchased bool      `json:"ifPurchased"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

type ItemDAO struct {
	db *sql.DB
}

func NewItemDAO(db *sql.DB) *ItemDAO {
	return &ItemDAO{db: db}
}

func (d *ItemDAO) GetItemsByCategory(category string, page, limit int) ([]*Item, error) {
	offset := (page - 1) * limit
	query := "SELECT id, title, price, explanation, image_url, uid, ifPurchased, category, created_at FROM items WHERE category = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := d.db.Query(query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Title, &item.Price, &item.Explanation, &item.ImageURL, &item.UID, &item.IfPurchased, &item.Category, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (d *ItemDAO) GetItemByID(id int) (*Item, error) {
	query := "SELECT id, title, price, explanation, image_url, uid, ifPurchased, category, created_at FROM items WHERE id = ?"
	row := d.db.QueryRow(query, id)

	var item Item
	err := row.Scan(&item.ID, &item.Title, &item.Price, &item.Explanation, &item.ImageURL, &item.UID, &item.IfPurchased, &item.Category, &item.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
