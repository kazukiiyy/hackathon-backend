# buyer_addressカラム追加手順

## 概要
`items`テーブルに`buyer_address`カラムを追加して、購入者のウォレットアドレスを保存できるようにします。これにより、ブロックチェーンからの購入イベントを正しく処理できるようになります。

## 追加されるカラム

### itemsテーブル
- `buyer_address` (VARCHAR(42), NULL): 購入者のウォレットアドレス
- インデックス: `idx_buyer_address` (検索パフォーマンス向上のため)

## SQL実行方法

### 推奨: 安全版スクリプトを使用（既存カラムをチェック）

```bash
cd hackathon-backend/scripts
mysql -h <HOST> -P <PORT> -u <USER> -p<PASSWORD> <DATABASE> < add_buyer_address_column_safe.sql
```

または

```bash
cd hackathon-backend/scripts
chmod +x run_sql_simple.sh
./run_sql_simple.sh add_buyer_address_column_safe.sql
```

### 方法1: 通常版SQLファイルを直接実行

```bash
cd hackathon-backend/scripts
mysql -h <HOST> -P <PORT> -u <USER> -p<PASSWORD> <DATABASE> < add_buyer_address_column.sql
```

### 方法2: スクリプトを使用

```bash
cd hackathon-backend/scripts
chmod +x run_sql_simple.sh
./run_sql_simple.sh add_buyer_address_column.sql
```

### 方法3: Cloud SQLに直接接続

Google Cloud SQLを使用している場合：

```bash
gcloud sql connect <INSTANCE_NAME> --user=<USER> --database=<DATABASE>
```

接続後、SQLファイルの内容をコピー&ペーストして実行してください。

## SQL内容

```sql
-- itemsテーブルにbuyer_addressカラムを追加
ALTER TABLE items 
ADD COLUMN buyer_address VARCHAR(42) NULL 
COMMENT '購入者のウォレットアドレス'
AFTER seller_address;

-- buyer_addressにインデックスを追加
CREATE INDEX idx_buyer_address ON items(buyer_address);
```

## 注意事項

- このSQLは既存のデータを保持します
- 既存のレコードの`buyer_address`は`NULL`になります
- 既に`buyer_address`カラムが存在する場合はエラーになりますが、無視して問題ありません
- インデックスが既に存在する場合もエラーになりますが、無視して問題ありません

## エラーが発生した場合

### エラー: "Duplicate column name 'buyer_address'"
- 既に`buyer_address`カラムが存在する場合は、このマイグレーションは不要です
- データベースのスキーマを確認してください

### エラー: "Duplicate key name 'idx_buyer_address'"
- 既にインデックスが存在する場合は、このマイグレーションは不要です
- データベースのスキーマを確認してください

## 実行後の確認

マイグレーション実行後、以下のSQLで確認できます：

```sql
-- buyer_addressカラムが存在するか確認
DESCRIBE items;

-- または
SHOW COLUMNS FROM items LIKE 'buyer_address';

-- インデックスが存在するか確認
SHOW INDEX FROM items WHERE Key_name = 'idx_buyer_address';
```

