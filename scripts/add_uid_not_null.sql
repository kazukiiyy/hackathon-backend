-- itemsテーブルのuidカラムをNOT NULL制約に変更
-- 既存のNULL値がある場合は、空文字列に更新してからNOT NULL制約を追加

-- 1. 既存のNULL値を空文字列に更新（安全のため）
UPDATE items 
SET uid = '' 
WHERE uid IS NULL;

-- 2. uidカラムにNOT NULL制約を追加
ALTER TABLE items 
MODIFY COLUMN uid VARCHAR(255) NOT NULL;

-- 3. インデックスを追加（オプション、パフォーマンス向上のため）
-- 既に存在する場合はエラーになるが、無視して問題ない
CREATE INDEX idx_items_uid ON items(uid);


