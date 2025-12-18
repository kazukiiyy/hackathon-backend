# SQL実行ガイド - uidカラムにNOT NULL制約を追加

## 概要

`items`テーブルの`uid`カラムにNOT NULL制約を追加します。
これにより、GET処理でuidがNULLの場合のエラーを防ぎます。

## 実行するSQL

```sql
-- 1. 既存のNULL値を空文字列に更新（安全のため）
UPDATE items 
SET uid = '' 
WHERE uid IS NULL;

-- 2. uidカラムにNOT NULL制約を追加
ALTER TABLE items 
MODIFY COLUMN uid VARCHAR(255) NOT NULL;

-- 3. インデックスを追加（オプション、パフォーマンス向上のため）
CREATE INDEX idx_items_uid ON items(uid);
```

## 実行方法

### 方法1: gcloud sql connect を使用（推奨・Cloud SQLの場合）

```bash
# Cloud SQLインスタンスに接続
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE>

# 接続後、SQLを実行
source add_uid_not_null.sql
```

または、SQLファイルを直接実行：

```bash
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE> < add_uid_not_null.sql
```

**実際の例：**
```bash
# インスタンス名を確認
gcloud sql instances list

# 接続（例）
gcloud sql connect my-instance --user=root --database=mydatabase

# 接続後、SQLを実行
source add_uid_not_null.sql
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
go run run_sql.go add_uid_not_null.sql
```

**実際の例：**
```bash
export MYSQL_USER="root"
export MYSQL_USER_PWD="your_password"
export MYSQL_DATABASE="mydatabase"
export INSTANCE_CONNECTION_NAME="term8-kazuki-tsukamoto:europe-west1:your-instance"

go run run_sql.go add_uid_not_null.sql
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
mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_USER_PWD" "$MYSQL_DATABASE" < add_uid_not_null.sql
```

**実際の例：**
```bash
# Cloud SQL Proxyが実行されている場合
export MYSQL_USER="root"
export MYSQL_USER_PWD="your_password"
export MYSQL_DATABASE="mydatabase"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"

mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_USER_PWD" "$MYSQL_DATABASE" < add_uid_not_null.sql
```

## 実行前の確認

既存のNULL値があるか確認：

```sql
-- uidがNULLのレコード数を確認
SELECT COUNT(*) FROM items WHERE uid IS NULL;

-- uidがNULLのレコードを確認（サンプル）
SELECT id, title, uid FROM items WHERE uid IS NULL LIMIT 10;
```

## 実行後の確認

カラムが正しく変更されたか確認：

```sql
-- itemsテーブルの構造を確認
DESCRIBE items;

-- uidカラムの詳細を確認
SHOW COLUMNS FROM items LIKE 'uid';

-- インデックスを確認
SHOW INDEX FROM items WHERE Key_name = 'idx_items_uid';

-- NULL値がないことを確認
SELECT COUNT(*) FROM items WHERE uid IS NULL;
```

## エラーが発生した場合

### カラムが既にNOT NULLの場合

```sql
-- カラムの制約を確認
SELECT 
    COLUMN_NAME,
    IS_NULLABLE,
    COLUMN_TYPE
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'items' 
AND COLUMN_NAME = 'uid';
```

既にNOT NULLの場合は、エラーが発生しますが、無視して問題ありません。

### インデックスが既に存在する場合

```sql
-- インデックスの存在を確認
SHOW INDEX FROM items WHERE Key_name = 'idx_items_uid';
```

既に存在する場合は、エラーが発生しますが、無視して問題ありません。

## 注意事項

- **既存のNULL値は空文字列に更新されます**
- 本番環境で実行する前に、必ずバックアップを取得してください
- 実行前にNULL値の数を確認し、必要に応じて適切なデフォルト値に変更してください
- `uid`にインデックスを追加することで、`GetItemsByUid`の検索パフォーマンスが向上します

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
GRANT ALTER, UPDATE, CREATE, INDEX ON your_database.* TO 'your_user'@'%';
```

### 大量のNULL値がある場合

大量のNULL値がある場合は、実行に時間がかかる可能性があります。
事前にNULL値の数を確認し、必要に応じて適切なデフォルト値に変更してください。

```sql
-- NULL値の数を確認
SELECT COUNT(*) FROM items WHERE uid IS NULL;

-- 必要に応じて、適切なデフォルト値に変更（例：'unknown'）
UPDATE items 
SET uid = 'unknown' 
WHERE uid IS NULL;
```




