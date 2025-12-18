-- itemsテーブルにbuyer_addressカラムを追加
-- 購入者のウォレットアドレスを保存するためのカラム

-- buyer_addressカラムを追加
-- 注意: 既に存在する場合はエラーになりますが、手動で確認してスキップしてください
ALTER TABLE items 
ADD COLUMN buyer_address VARCHAR(42) NULL 
COMMENT '購入者のウォレットアドレス'
AFTER seller_address;

-- buyer_addressにインデックスを追加（検索パフォーマンス向上のため）
-- 注意: 既に存在する場合はエラーになりますが、手動で確認してスキップしてください
CREATE INDEX idx_buyer_address ON items(buyer_address);

