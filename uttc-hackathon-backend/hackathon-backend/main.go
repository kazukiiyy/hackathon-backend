package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	getItemDao "uttc-hackathon-backend/dao/getItems"
	likesDao "uttc-hackathon-backend/dao/likes"
	messagesDao "uttc-hackathon-backend/dao/messages"
	postItemsDao "uttc-hackathon-backend/dao/postItems"
	postUserDao "uttc-hackathon-backend/dao/postUser"
	purchaseItemDao "uttc-hackathon-backend/dao/purchaseItem"
	getItemHdr "uttc-hackathon-backend/handlers/getItems"
	likesHdr "uttc-hackathon-backend/handlers/likes"
	messagesHdr "uttc-hackathon-backend/handlers/messages"
	postItemsHdr "uttc-hackathon-backend/handlers/postItems"
	postUserHdr "uttc-hackathon-backend/handlers/postUser"
	purchaseItemHdr "uttc-hackathon-backend/handlers/purchaseItem"
	blockchainHdr "uttc-hackathon-backend/handlers/blockchain"
	blockchainUc "uttc-hackathon-backend/usecase/blockchain"
	geminiHdr "uttc-hackathon-backend/handlers/gemini"
	geminiUc "uttc-hackathon-backend/usecase/gemini"
	getItemUc "uttc-hackathon-backend/usecase/getItems"
	likesUc "uttc-hackathon-backend/usecase/likes"
	messagesUc "uttc-hackathon-backend/usecase/messages"
	postItemsUc "uttc-hackathon-backend/usecase/postItems"
	postUserUc "uttc-hackathon-backend/usecase/postUser"
	purchaseItemUc "uttc-hackathon-backend/usecase/purchaseItem"

	_ "github.com/go-sql-driver/mysql"
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

	// 既存のハンドラー設定
	itemDAO := postItemsDao.NewItemDAO(db)
	itemUsecase := postItemsUc.NewItemUsecase(itemDAO)
	itemHandler := postItemsHdr.NewItemHandler(itemUsecase)

	getItemDAO := getItemDao.NewItemDAO(db)
	getItemUsecase := getItemUc.NewItemUsecase(getItemDAO)
	getItemHandler := getItemHdr.NewItemHandler(getItemUsecase)

	userDAO := postUserDao.NewUserDAO(db)
	userUsecase := postUserUc.NewUserUsecase(userDAO)
	userHandler := postUserHdr.NewUserHandler(userUsecase)

	purchaseDAO := purchaseItemDao.NewPurchaseDAO(db)
	purchaseUsecase := purchaseItemUc.NewPurchaseUsecase(purchaseDAO)
	purchaseHandler := purchaseItemHdr.NewPurchaseHandler(purchaseUsecase)

	messageDAO := messagesDao.NewMessageDAO(db)
	messageUsecase := messagesUc.NewMessageUsecase(messageDAO)
	messageHandler := messagesHdr.NewMessageHandler(messageUsecase)

	likeDAO := likesDao.NewLikeDAO(db)
	likeUsecase := likesUc.NewLikeUsecase(likeDAO)
	likeHandler := likesHdr.NewLikeHandler(likeUsecase)

	// Blockchain handler
	blockchainUsecase := blockchainUc.NewBlockchainUsecase(itemDAO, purchaseDAO)
	blockchainHandler := blockchainHdr.NewBlockchainHandler(blockchainUsecase)

	// Gemini handler
	geminiUsecase := geminiUc.NewGeminiUsecase()
	geminiHandler := geminiHdr.NewGeminiHandler(geminiUsecase)

	// HTTPルーティング
	http.HandleFunc("/postItems", itemHandler.CreateItem)
	http.HandleFunc("/uploadImage", itemHandler.UploadImage)
	http.HandleFunc("/getItems", getItemHandler.GetItems)
	http.HandleFunc("/getItems/latest", getItemHandler.GetLatestItems)
	http.HandleFunc("/getItems/", getItemHandler.GetItemByID)
	http.HandleFunc("/register", userHandler.RegisterUser)
	http.HandleFunc("/items/", purchaseHandler.PurchaseItem)
	http.HandleFunc("/purchases", purchaseHandler.GetPurchasedItems)
	http.HandleFunc("/messages", messageHandler.GetMessages)
	http.HandleFunc("/messages/send", messageHandler.SendMessage)
	http.HandleFunc("/messages/read", messageHandler.MarkAsRead)
	http.HandleFunc("/messages/conversations", messageHandler.GetConversations)
	http.HandleFunc("/likes", likeHandler.HandleLike)
	http.HandleFunc("/likes/status", likeHandler.GetLikeStatus)
	http.HandleFunc("/likes/user", likeHandler.GetUserLikes)
	// Blockchain endpoints
	http.HandleFunc("/api/v1/blockchain/item-listed", blockchainHandler.HandleItemListed)
	http.HandleFunc("/api/v1/blockchain/item-purchased", blockchainHandler.HandleItemPurchased)
	// Gemini endpoint
	http.HandleFunc("/api/v1/gemini/generate", geminiHandler.GenerateContent)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	standardRouter := http.DefaultServeMux
	
	loggingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		corsMiddleware(standardRouter).ServeHTTP(lrw, r)
		
		if lrw.statusCode >= 400 {
			log.Printf("Request failed: %s %s status=%d", r.Method, r.URL.Path, lrw.statusCode)
		}
	})
	
	finalHandler := loggingHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)

	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatal(err)
	}
}
