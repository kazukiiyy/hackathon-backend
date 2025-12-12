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

func (d *ItemDAO) InsertItem(title string, price int, explanation string, imagePath string, uid string, ifPurchased bool, category string) error {
	query := "INSERT INTO items (title, price, explanation, image_url, uid, ifPurchased, category) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := d.db.Exec(query, title, price, explanation, imagePath, uid, ifPurchased, category)
	return err
}
