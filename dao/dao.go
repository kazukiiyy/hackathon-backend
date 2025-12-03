package dao

import (
	"database/sql"
)

type UserDAO struct {
	db *sql.DB
}

func NewItemDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) InsertItem(title string, price int, explanation string, imagePath string) error {
	query := "INSERT INTO items (title, price, explanation, image_url) VALUES (?, ?, ?, ?)"
	_, err := d.db.Exec(query, title, price, explanation, imagePath)
	return err
}
