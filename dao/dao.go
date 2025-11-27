package dao

import (
	"database/sql"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(dbMain *sql.DB) *UserDAO {
	return &UserDAO{db: dbMain}
}
