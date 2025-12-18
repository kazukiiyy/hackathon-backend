# GCPコンソールからSQLを実行する手順

## 手順

1. **Google Cloud Consoleにアクセス**
   - https://console.cloud.google.com/ にアクセス
   - プロジェクトを選択

2. **Cloud SQLインスタンスに移動**
   - 左側のメニューから「SQL」を選択
   - または検索バーで「Cloud SQL」を検索

3. **インスタンスを選択**
   - 使用しているCloud SQLインスタンスをクリック

4. **データベースを選択**
   - 上部のタブから「データベース」を選択
   - 使用しているデータベース（例: `mydatabase`）をクリック

5. **クエリエディタを開く**
   - データベースの詳細ページで「クエリ」タブをクリック
   - または、データベース名の横にある「クエリ」ボタンをクリック

6. **SQLを実行**
   - 以下のSQL文をコピー＆ペースト
   - 「実行」ボタンをクリック

## 実行するSQL

```sql
-- 全テーブルのデータを削除（外部キー制約を考慮して順序を守る）

-- 外部キー制約を一時的に無効化
SET FOREIGN_KEY_CHECKS = 0;

-- データを削除（テーブル構造は保持）
TRUNCATE TABLE item_images;
TRUNCATE TABLE purchases;
TRUNCATE TABLE likes;
TRUNCATE TABLE messages;
TRUNCATE TABLE items;
TRUNCATE TABLE users;

-- 外部キー制約を再有効化
SET FOREIGN_KEY_CHECKS = 1;
```

## 注意事項

- この操作は**すべてのデータを削除**します
- テーブル構造は保持されます
- 実行前にバックアップを取ることを推奨します


