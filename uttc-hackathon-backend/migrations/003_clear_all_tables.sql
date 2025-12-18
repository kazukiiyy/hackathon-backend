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
