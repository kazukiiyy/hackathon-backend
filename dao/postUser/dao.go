package postUser

import (
	"database/sql"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) InsertUser(uid string, nickname string, sex string, birthyear int, birthdate int) error {
	query := "INSERT INTO users (uid, nickname, sex, birthyear, birthdate) VALUES (?, ?, ?, ?, ?)"
	_, err := d.db.Exec(query, uid, nickname, sex, birthyear, birthdate)
	return err
}
