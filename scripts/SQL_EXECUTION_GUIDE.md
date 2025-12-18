# SQL実行ガイド - ブロックチェーン関連カラム追加

## スマートコントラクトの確認

✅ **確認済み**: `ItemListed`イベントには以下のフィールドが含まれています：
- `itemId` (chain_item_id)
- `tokenId` (token_id)
- `seller` (seller_address)
- `title`, `price`, `explanation`, `imageUrl`, `uid`, `category`, `createdAt`

## 実行するSQL

```sql
-- itemsテーブルにブロックチェーン関連のカラムを追加
ALTER TABLE items 
ADD COLUMN chain_item_id BIGINT NULL 
AFTER status;

ALTER TABLE items 
ADD COLUMN seller_address VARCHAR(255) NULL 
AFTER chain_item_id;

ALTER TABLE items 
ADD COLUMN token_id BIGINT NULL 
AFTER seller_address;

-- chain_item_idにインデックスを追加
CREATE INDEX idx_items_chain_item_id ON items(chain_item_id);

-- item_imagesテーブルにchain_item_idカラムを追加（オプション）
ALTER TABLE item_images 
ADD COLUMN chain_item_id BIGINT NULL 
AFTER item_id;
```

## 実行方法

### 方法1: gcloud sql connect を使用（推奨・Cloud SQLの場合）

```bash
# Cloud SQLインスタンスに接続
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE>

# 接続後、SQLを実行
# 以下のSQLをコピー&ペーストして実行してください
```

または、SQLファイルを直接実行：

```bash
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE> < add_chain_columns.sql
```

**実際の例：**
```bash
# インスタンス名を確認
gcloud sql instances list

# 接続（例）
gcloud sql connect my-instance --user=root --database=mydatabase

# 接続後、SQLを実行
source add_chain_columns.sql
# または
mysql> ALTER TABLE items ADD COLUMN chain_item_id BIGINT NULL AFTER status;
mysql> ALTER TABLE items ADD COLUMN seller_address VARCHAR(255) NULL AFTER chain_item_id;
mysql> ALTER TABLE items ADD COLUMN token_id BIGINT NULL AFTER seller_address;
mysql> CREATE INDEX idx_items_chain_item_id ON items(chain_item_id);
mysql> ALTER TABLE item_images ADD COLUMN chain_item_id BIGINT NULL AFTER item_id;
```

### 方法2: Goスクリプトを使用（Cloud SQL Unix Socket経由）

```bash
cd hackathon-backend/scripts

# 環境変数を設定
export MYSQL_USER="your_user"
export MYSQL_USER_PWD="your_password"
export MYSQL_DATABASE="your_database"
export INSTANCE_CONNECTION_NAME="project:region:instance"

# SQLファイルを実行
go run run_sql.go add_chain_columns.sql
```

**実際の例：**
```bash
export MYSQL_USER="root"
export MYSQL_USER_PWD="your_password"
export MYSQL_DATABASE="mydatabase"
export INSTANCE_CONNECTION_NAME="term8-kazuki-tsukamoto:europe-west1:your-instance"

go run run_sql.go add_chain_columns.sql
```

### 方法3: MySQLクライアントで直接実行

```bash
cd hackathon-backend/scripts

# 環境変数を設定
export MYSQL_USER="your_user"
export MYSQL_USER_PWD="your_password"
export MYSQL_DATABASE="your_database"
export MYSQL_HOST="127.0.0.1"  # Cloud SQL Proxy経由の場合
export MYSQL_PORT="3306"

# SQLファイルを実行
mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_USER_PWD" "$MYSQL_DATABASE" < add_chain_columns.sql
```

**実際の例：**
```bash
# Cloud SQL Proxyが実行されている場合
export MYSQL_USER="root"
export MYSQL_USER_PWD="your_password"
export MYSQL_DATABASE="mydatabase"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"

mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_USER_PWD" "$MYSQL_DATABASE" < add_chain_columns.sql
```

### 方法4: 実行スクリプトを使用

```bash
cd hackathon-backend/scripts
chmod +x execute_add_chain_columns.sh
./execute_add_chain_columns.sh
```

## 実行後の確認

カラムが正しく追加されたか確認：

```sql
-- itemsテーブルの構造を確認
DESCRIBE items;

-- 特定のカラムを確認
SHOW COLUMNS FROM items LIKE 'chain_item_id';
SHOW COLUMNS FROM items LIKE 'seller_address';
SHOW COLUMNS FROM items LIKE 'token_id';

-- インデックスを確認
SHOW INDEX FROM items WHERE Key_name = 'idx_items_chain_item_id';
```

## エラーが発生した場合

### カラムが既に存在する場合

```sql
-- カラムの存在を確認
SELECT COLUMN_NAME 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'items' 
AND COLUMN_NAME IN ('chain_item_id', 'seller_address', 'token_id');

-- 既に存在する場合は、そのカラムをスキップして他のカラムのみ追加
```

### item_imagesテーブルのchain_item_idが既に存在する場合

エラーが発生しますが、無視して問題ありません。以下のSQLで確認できます：

```sql
-- item_imagesテーブルにchain_item_idが存在するか確認
SELECT COLUMN_NAME 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'item_images' 
AND COLUMN_NAME = 'chain_item_id';
```

## 注意事項

- 既存のデータは保持されます
- 既存レコードの新規カラムは`NULL`になります
- `chain_item_id`にインデックスを追加することで、検索パフォーマンスが向上します
- 本番環境で実行する前に、必ずバックアップを取得してください

## トラブルシューティング

### 接続エラーが発生する場合

1. Cloud SQL Proxyが実行されているか確認
2. 環境変数が正しく設定されているか確認
3. インスタンス名が正しいか確認

### 権限エラーが発生する場合

```sql
-- 必要な権限を確認
SHOW GRANTS FOR 'your_user'@'%';

-- 権限がない場合は、管理者に依頼
GRANT ALTER, CREATE, INDEX ON your_database.* TO 'your_user'@'%';
```


