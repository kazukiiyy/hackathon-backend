-- itemsテーブルにブロックチェーン関連のカラムを追加
-- chain_item_id: ブロックチェーン上の商品ID
-- seller_address: 売り手のウォレットアドレス
-- token_id: NFTのトークンID

-- chain_item_idカラムを追加
ALTER TABLE items 
ADD COLUMN chain_item_id BIGINT NULL 
AFTER status;

-- seller_addressカラムを追加
ALTER TABLE items 
ADD COLUMN seller_address VARCHAR(255) NULL 
AFTER chain_item_id;

-- token_idカラムを追加
ALTER TABLE items 
ADD COLUMN token_id BIGINT NULL 
AFTER seller_address;

-- chain_item_idにインデックスを追加（検索パフォーマンス向上のため）
CREATE INDEX idx_items_chain_item_id ON items(chain_item_id);

-- item_imagesテーブルにchain_item_idカラムを追加（オプション）
-- 既に存在する場合はエラーになるが、無視して問題ない
ALTER TABLE item_images 
ADD COLUMN chain_item_id BIGINT NULL 
AFTER item_id;

