#!/bin/bash

# 簡単なSQL実行スクリプト
# 使用方法: ./run_sql_simple.sh <sql_file_path>

if [ $# -lt 1 ]; then
    echo "Usage: $0 <sql_file_path>"
    exit 1
fi

SQL_FILE=$1

# 環境変数を確認
if [ -z "$MYSQL_USER" ] || [ -z "$MYSQL_USER_PWD" ] || [ -z "$MYSQL_DATABASE" ]; then
    echo "Error: Required environment variables not set"
    echo "Please set: MYSQL_USER, MYSQL_USER_PWD, MYSQL_DATABASE"
    exit 1
fi

# MySQL接続（Cloud SQL Proxy経由の場合）
# 注意: Cloud SQL Proxyが実行されている必要があります
# または、パブリックIPを使用する場合は接続文字列を変更してください

MYSQL_HOST=${MYSQL_HOST:-"127.0.0.1"}
MYSQL_PORT=${MYSQL_PORT:-"3306"}

echo "Connecting to MySQL at $MYSQL_HOST:$MYSQL_PORT..."
echo "Executing SQL file: $SQL_FILE"

mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_USER_PWD" "$MYSQL_DATABASE" < "$SQL_FILE"

if [ $? -eq 0 ]; then
    echo "SQL file executed successfully!"
else
    echo "Error executing SQL file"
    exit 1
fi


