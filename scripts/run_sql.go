package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run run_sql.go <sql_file_path>")
	}

	sqlFile := os.Args[1]

	// 環境変数から接続情報を取得
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_USER_PWD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")

	if mysqlUser == "" || mysqlUserPwd == "" || mysqlDatabase == "" || instanceConnectionName == "" {
		log.Fatal("Required environment variables: MYSQL_USER, MYSQL_USER_PWD, MYSQL_DATABASE, INSTANCE_CONNECTION_NAME")
	}

	// SQLファイルを読み込む
	sqlContent, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// DSNを構築（Cloud SQL用）
	dsn := fmt.Sprintf(
		"%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
		mysqlUser,
		mysqlUserPwd,
		instanceConnectionName,
		mysqlDatabase,
	)

	// データベースに接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 接続をテスト
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully!")

	// SQL文を分割（セミコロンで区切る）
	statements := strings.Split(string(sqlContent), ";")

	// 各SQL文を実行
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		log.Printf("Executing statement %d/%d...", i+1, len(statements))
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Error executing statement %d: %v", i+1, err)
			log.Printf("Statement: %s", stmt)
			log.Fatal("Aborting...")
		}
	}

	log.Println("All SQL statements executed successfully!")
}


