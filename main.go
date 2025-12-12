package main

import (
	"database/sql"
	"fmt"
	"net/http"
	getItemDao "uttc-hackathon-backend/dao/getItems"
	postItemsDao "uttc-hackathon-backend/dao/postItems"
	postUserDao "uttc-hackathon-backend/dao/postUser"
	getItemHdr "uttc-hackathon-backend/handlers/getItems"
	postItemsHdr "uttc-hackathon-backend/handlers/postItems"
	postUserHdr "uttc-hackathon-backend/handlers/postUser"
	getItemUc "uttc-hackathon-backend/usecase/getItems"
	postItemsUc "uttc-hackathon-backend/usecase/postItems"
	postUserUc "uttc-hackathon-backend/usecase/postUser"

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

	getItemDAO := getItemDao.NewItemDAO(db)
	getItemUsecase := getItemUc.NewItemUsecase(getItemDAO)
	getItemHandler := getItemHdr.NewItemHandler(getItemUsecase)

	userDAO := postUserDao.NewUserDAO(db)
	userUsecase := postUserUc.NewUserUsecase(userDAO)
	userHandler := postUserHdr.NewUserHandler(userUsecase)

	http.HandleFunc("/postItems", itemHandler.CreateItem)
	http.HandleFunc("/getItems", getItemHandler.GetItems)
	http.HandleFunc("/register", userHandler.RegisterUser)
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
