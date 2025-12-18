-- purchasesテーブルにbuyer_addressカラムを追加（シンプル版）
-- GCPコンソールのクエリエディタで直接実行可能

-- buyer_addressカラムを追加（既に存在する場合はエラーになりますが、無視して問題ありません）
ALTER TABLE purchases 
ADD COLUMN buyer_address VARCHAR(42) NULL 
COMMENT '購入者ウォレットアドレス'
AFTER buyer_uid;



