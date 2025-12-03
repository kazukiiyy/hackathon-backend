package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"uttc-hackathon-backend/dao"
	handler "uttc-hackathon-backend/handlers"

	_ "github.com/go-sql-driver/mysql"

	"log"
	"os"
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

	itemDAO := dao.NewItemDAO(_db)
	itemHandler := handler.ItemHandler(itemDAO)

	http.HandleFunc("/items", itemHandler.CreateItem)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
