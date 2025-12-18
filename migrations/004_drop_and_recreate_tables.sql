-- 全テーブルを削除して再作成（完全リセット用）

-- 外部キー制約を一時的に無効化
SET FOREIGN_KEY_CHECKS = 0;

-- テーブルを削除（存在する場合）
DROP TABLE IF EXISTS item_images;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;

-- 外部キー制約を再有効化
SET FOREIGN_KEY_CHECKS = 1;

-- itemsテーブル
CREATE TABLE items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    chain_item_id BIGINT UNIQUE COMMENT 'スマートコントラクト上の商品ID',
    token_id BIGINT COMMENT 'NFTのトークンID',
    title VARCHAR(255) NOT NULL COMMENT '商品タイトル',
    price VARCHAR(78) NOT NULL COMMENT '価格（Wei単位）',
    explanation TEXT COMMENT '商品説明',
    uid VARCHAR(255) NOT NULL COMMENT 'ユーザーID（Firebase等）',
    seller_address VARCHAR(42) COMMENT '出品者のウォレットアドレス',
    buyer_address VARCHAR(42) COMMENT '購入者のウォレットアドレス',
    status ENUM('listed', 'purchased', 'completed', 'cancelled') DEFAULT 'listed' COMMENT '商品の状態',
    category VARCHAR(100) COMMENT 'カテゴリー',
    like_count INT DEFAULT 0 COMMENT 'いいね数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    INDEX idx_chain_item_id (chain_item_id),
    INDEX idx_status (status),
    INDEX idx_seller_address (seller_address),
    INDEX idx_buyer_address (buyer_address),
    INDEX idx_category (category),
    INDEX idx_uid (uid),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- item_imagesテーブル
CREATE TABLE item_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT COMMENT '旧形式の商品ID',
    chain_item_id BIGINT COMMENT 'スマートコントラクト上の商品ID',
    image_url VARCHAR(500) NOT NULL COMMENT '画像URL',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_item_id (item_id),
    INDEX idx_chain_item_id (chain_item_id),
    FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- usersテーブル
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uid VARCHAR(255) UNIQUE NOT NULL COMMENT 'Firebase UID',
    nickname VARCHAR(100) COMMENT 'ニックネーム',
    wallet_address VARCHAR(42) COMMENT 'ウォレットアドレス',
    sex VARCHAR(10) COMMENT '性別',
    birthyear INT COMMENT '生年',
    birthdate VARCHAR(10) COMMENT '誕生日',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uid (uid),
    INDEX idx_wallet_address (wallet_address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- purchasesテーブル
CREATE TABLE purchases (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT COMMENT '旧形式の商品ID',
    chain_item_id BIGINT COMMENT 'スマートコントラクト上の商品ID',
    buyer_uid VARCHAR(255) NOT NULL COMMENT '購入者UID',
    buyer_address VARCHAR(42) COMMENT '購入者ウォレットアドレス',
    purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_item_id (item_id),
    INDEX idx_chain_item_id (chain_item_id),
    INDEX idx_buyer_uid (buyer_uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- messagesテーブル
CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_uid VARCHAR(255) NOT NULL,
    receiver_uid VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sender_uid (sender_uid),
    INDEX idx_receiver_uid (receiver_uid),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- likesテーブル
CREATE TABLE likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT COMMENT '旧形式の商品ID',
    chain_item_id BIGINT COMMENT 'スマートコントラクト上の商品ID',
    uid VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_like (item_id, uid),
    INDEX idx_item_id (item_id),
    INDEX idx_chain_item_id (chain_item_id),
    INDEX idx_uid (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


