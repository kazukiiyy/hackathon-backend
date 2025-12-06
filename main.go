package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"uttc-hackathon-backend/dao"
	handler "uttc-hackathon-backend/handlers"
	"uttc-hackathon-backend/usecase/post_item"

	_ "github.com/go-sql-driver/mysql"

	"log"
	"os"
)

func main() {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_USER_PWD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlHost := os.Getenv("MYSQL_HOST")

	_db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", mysqlUser, mysqlUserPwd, mysqlHost, mysqlDatabase))
	if err != nil {
		log.Fatalf("fail:sql.Open,%v/n", err)
	}

	defer _db.Close()

	if err := _db.Ping(); err != nil {
		log.Fatalf("fail:_db.Ping,%v/n", err)
	}
	log.Println("Connected to mysql")

	itemDAO := dao.NewItemDAO(_db)
	itemUsecase := post_item.NewItemUsecase(itemDAO)
	itemHandler := handler.NewItemHandler(itemUsecase)

	http.HandleFunc("/items", itemHandler.CreateItem)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
