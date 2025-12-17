-- ブロックチェーン連携用のカラムを追加するマイグレーション

-- itemsテーブルにブロックチェーン関連カラムを追加
ALTER TABLE items
    ADD COLUMN chain_item_id BIGINT UNIQUE COMMENT 'スマートコントラクト上の商品ID',
    ADD COLUMN token_id BIGINT COMMENT 'NFTのトークンID',
    ADD COLUMN seller_address VARCHAR(42) COMMENT '出品者のウォレットアドレス',
    ADD COLUMN buyer_address VARCHAR(42) COMMENT '購入者のウォレットアドレス',
    ADD COLUMN status ENUM('listed', 'purchased', 'completed', 'cancelled') DEFAULT 'listed' COMMENT '商品の状態',
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時';

-- priceカラムをVARCHARに変更（Wei単位の大きな数値を扱うため）
-- 既存データがある場合は先にバックアップを取ること
ALTER TABLE items MODIFY COLUMN price VARCHAR(78) COMMENT '価格（Wei単位）';

-- item_imagesテーブルにchain_item_idカラムを追加
ALTER TABLE item_images
    ADD COLUMN chain_item_id BIGINT COMMENT 'スマートコントラクト上の商品ID',
    ADD INDEX idx_chain_item_id (chain_item_id);

-- インデックスの追加
ALTER TABLE items
    ADD INDEX idx_chain_item_id (chain_item_id),
    ADD INDEX idx_status (status),
    ADD INDEX idx_seller_address (seller_address),
    ADD INDEX idx_buyer_address (buyer_address);
