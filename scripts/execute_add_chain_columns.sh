#!/bin/bash

# ブロックチェーン関連カラム追加スクリプト
# 使用方法: ./execute_add_chain_columns.sh

set -e  # エラーが発生したら終了

echo "=========================================="
echo "ブロックチェーン関連カラム追加スクリプト"
echo "=========================================="

# SQLファイルのパス
SQL_FILE="$(dirname "$0")/add_chain_columns.sql"

if [ ! -f "$SQL_FILE" ]; then
    echo "エラー: SQLファイルが見つかりません: $SQL_FILE"
    exit 1
fi

# 方法1: gcloud sql connect を使用（推奨）
if command -v gcloud &> /dev/null; then
    echo ""
    echo "方法1: gcloud sql connect を使用"
    echo "----------------------------------------"
    echo "以下のコマンドを実行してください:"
    echo ""
    echo "gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE>"
    echo ""
    echo "接続後、以下のSQLを実行してください:"
    echo ""
    cat "$SQL_FILE"
    echo ""
    echo "または、以下のコマンドで直接実行:"
    echo ""
    echo "gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE> < $SQL_FILE"
    echo ""
fi

# 方法2: 環境変数を使用してMySQLクライアントで実行
if [ -n "$MYSQL_USER" ] && [ -n "$MYSQL_USER_PWD" ] && [ -n "$MYSQL_DATABASE" ]; then
    echo ""
    echo "方法2: MySQLクライアントで実行"
    echo "----------------------------------------"
    
    MYSQL_HOST=${MYSQL_HOST:-"127.0.0.1"}
    MYSQL_PORT=${MYSQL_PORT:-"3306"}
    
    echo "接続先: $MYSQL_HOST:$MYSQL_PORT"
    echo "データベース: $MYSQL_DATABASE"
    echo "ユーザー: $MYSQL_USER"
    echo ""
    echo "SQLファイルを実行します..."
    
    mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_USER_PWD" "$MYSQL_DATABASE" < "$SQL_FILE"
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✅ SQLファイルが正常に実行されました！"
        echo ""
        echo "追加されたカラムを確認するには:"
        echo "  DESCRIBE items;"
        echo "  または"
        echo "  SHOW COLUMNS FROM items LIKE 'chain_item_id';"
    else
        echo ""
        echo "❌ SQLファイルの実行に失敗しました"
        exit 1
    fi
else
    echo ""
    echo "方法2: MySQLクライアントで実行（環境変数が必要）"
    echo "----------------------------------------"
    echo "以下の環境変数を設定してください:"
    echo "  export MYSQL_USER=<ユーザー名>"
    echo "  export MYSQL_USER_PWD=<パスワード>"
    echo "  export MYSQL_DATABASE=<データベース名>"
    echo "  export MYSQL_HOST=<ホスト> (オプション、デフォルト: 127.0.0.1)"
    echo "  export MYSQL_PORT=<ポート> (オプション、デフォルト: 3306)"
    echo ""
fi

# 方法3: Goスクリプトを使用
if command -v go &> /dev/null; then
    echo ""
    echo "方法3: Goスクリプトを使用"
    echo "----------------------------------------"
    echo "以下の環境変数を設定して実行:"
    echo "  export MYSQL_USER=<ユーザー名>"
    echo "  export MYSQL_USER_PWD=<パスワード>"
    echo "  export MYSQL_DATABASE=<データベース名>"
    echo "  export INSTANCE_CONNECTION_NAME=<接続名>"
    echo ""
    echo "  cd $(dirname "$0")"
    echo "  go run run_sql.go add_chain_columns.sql"
    echo ""
fi

echo ""
echo "=========================================="
echo "実行完了"
echo "=========================================="

