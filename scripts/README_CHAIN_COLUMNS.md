# ブロックチェーン関連カラム追加手順

## 概要
`items`テーブルにブロックチェーン関連のカラムを追加して、スマートコントラクトからのイベントを処理できるようにします。

## 追加されるカラム

### itemsテーブル
- `chain_item_id` (BIGINT, NULL): ブロックチェーン上の商品ID
- `seller_address` (VARCHAR(255), NULL): 売り手のウォレットアドレス
- `token_id` (BIGINT, NULL): NFTのトークンID

### item_imagesテーブル（オプション）
- `chain_item_id` (BIGINT, NULL): ブロックチェーン上の商品ID（画像との関連付け用）

## SQL実行方法

### 方法1: SQLファイルを直接実行

```bash
cd hackathon-backend/scripts
mysql -h <HOST> -P <PORT> -u <USER> -p<PASSWORD> <DATABASE> < add_chain_columns.sql
```

### 方法2: スクリプトを使用

```bash
cd hackathon-backend/scripts
chmod +x run_sql_simple.sh
./run_sql_simple.sh add_chain_columns.sql
```

### 方法3: Cloud SQLに直接接続

Google Cloud SQLを使用している場合：

```bash
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE>
```

接続後、SQLファイルの内容をコピー&ペーストして実行してください。

## SQL内容

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

## 注意事項

- このSQLは既存のデータを保持します
- 既存のレコードの`chain_item_id`、`seller_address`、`token_id`は`NULL`になります
- `chain_item_id`にインデックスを追加することで、検索パフォーマンスが向上します
- `item_images`テーブルの`chain_item_id`カラムは既に存在する場合、エラーになりますが無視して問題ありません

## 実行後の確認

以下のクエリでカラムが正しく追加されたか確認できます：

```sql
-- itemsテーブルの構造を確認
DESCRIBE items;

-- chain_item_idカラムが存在するか確認
SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'items' 
AND COLUMN_NAME IN ('chain_item_id', 'seller_address', 'token_id');
```


