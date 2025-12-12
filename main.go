package main

import (
	"database/sql"
	"fmt"
	"net/http"
	postItemsDao "uttc-hackathon-backend/dao/postItems"
	postItemsHdr "uttc-hackathon-backend/handlers/postItems"
	postItemsUc "uttc-hackathon-backend/usecase/postItems"

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

	itemDAO := postItemsDao.NewItemDAO(db)
	itemUsecase := postItemsUc.NewItemUsecase(itemDAO)
	itemHandler := postItemsHdr.NewItemHandler(itemUsecase)

	http.HandleFunc("/items", itemHandler.CreateItem)
	standardRouter := http.DefaultServeMux
	finalHandler := corsMiddleware(standardRouter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル用 fallback
	}

	log.Printf("Server listening on port %s", port)

	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatal(err)
	}

}
