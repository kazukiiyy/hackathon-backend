-- itemsテーブルにbuyer_addressカラムを追加（安全版）
-- 購入者のウォレットアドレスを保存するためのカラム
-- 既に存在する場合はスキップします

-- buyer_addressカラムが存在するか確認してから追加
SET @dbname = DATABASE();
SET @tablename = 'items';
SET @columnname = 'buyer_address';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (COLUMN_NAME = @columnname)
  ) > 0,
  'SELECT "Column buyer_address already exists in items table" AS result;',
  CONCAT('ALTER TABLE ', @tablename, ' ADD COLUMN ', @columnname, ' VARCHAR(42) NULL COMMENT ''購入者のウォレットアドレス'' AFTER seller_address;')
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- buyer_addressインデックスが存在するか確認してから追加
SET @indexname = 'idx_buyer_address';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (INDEX_NAME = @indexname)
  ) > 0,
  'SELECT "Index idx_buyer_address already exists on items table" AS result;',
  CONCAT('CREATE INDEX ', @indexname, ' ON ', @tablename, '(', @columnname, ');')
));
PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

