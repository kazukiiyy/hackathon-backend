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
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
		mysqlUser,
		mysqlUserPwd,
		instanceConnectionName,
		mysqlDatabase,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}
	defer db.Close()

	log.Println("Connected to Cloud SQL!")

	itemDAO := dao.NewItemDAO(db)
	itemUsecase := post_item.NewItemUsecase(itemDAO)
	itemHandler := handler.NewItemHandler(itemUsecase)

	http.HandleFunc("/items", itemHandler.CreateItem)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル用 fallback
	}

	log.Printf("Server listening on port %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
