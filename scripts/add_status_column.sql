-- itemsテーブルにstatusカラムを追加
-- statusカラムは商品の状態を表す（listed, purchased, completed, cancelled）

ALTER TABLE items 
ADD COLUMN status VARCHAR(50) DEFAULT 'listed' 
AFTER category;

-- 既存のレコードにデフォルト値を設定
UPDATE items 
SET status = 'listed' 
WHERE status IS NULL;

-- インデックスを追加（オプション、パフォーマンス向上のため）
CREATE INDEX idx_items_status ON items(status);

