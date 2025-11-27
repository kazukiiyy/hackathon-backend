package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"log"
	"os"
	"uttc-hackathon-backend/dao"
)

func main() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_USER_PWD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(localhost:3306)/%s", mysqlUser, mysqlUserPwd, mysqlDatabase))
	if err != nil {
		log.Fatalf("fail:sql.Open,%v/n", err)
	}

	defer _db.Close()

	if err := _db.Ping(); err != nil {
		log.Fatalf("fail:_db.Ping,%v/n", err)
	}
	log.Println("Connected to mysql")

	userDAO := dao.NewUserDAO(_db)

	_ = userDAO

}
