# statusカラム追加手順

## 概要
`items`テーブルに`status`カラムを追加して、商品の状態（listed, purchased, completed, cancelled）を管理できるようにします。

## SQL実行方法

### 方法1: SQLファイルを直接実行

```bash
cd hackathon-backend/scripts
mysql -h <HOST> -P <PORT> -u <USER> -p<PASSWORD> <DATABASE> < add_status_column.sql
```

### 方法2: スクリプトを使用

```bash
cd hackathon-backend/scripts
chmod +x run_sql_simple.sh
./run_sql_simple.sh add_status_column.sql
```

### 方法3: Cloud SQLに直接接続

Google Cloud SQLを使用している場合：

```bash
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE>
```

接続後、SQLファイルの内容をコピー&ペーストして実行してください。

## SQL内容

```sql
-- itemsテーブルにstatusカラムを追加
ALTER TABLE items 
ADD COLUMN status VARCHAR(50) DEFAULT 'listed' 
AFTER category;

-- 既存のレコードにデフォルト値を設定
UPDATE items 
SET status = 'listed' 
WHERE status IS NULL;

-- インデックスを追加（オプション）
CREATE INDEX idx_items_status ON items(status);
```

## 注意事項

- このSQLは既存のデータを保持します
- 既存のレコードには`status = 'listed'`が設定されます
- `status`カラムは`VARCHAR(50)`型で、デフォルト値は`'listed'`です

