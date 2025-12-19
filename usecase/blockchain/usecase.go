package blockchain

import (
	"fmt"
	"log"
	"math/big"
	postItemsDao "uttc-hackathon-backend/dao/postItems"
	purchaseItemDao "uttc-hackathon-backend/dao/purchaseItem"
)

type BlockchainUsecase struct {
	itemDAO     *postItemsDao.ItemDAO
	purchaseDAO *purchaseItemDao.PurchaseDAO
}

func NewBlockchainUsecase(itemDAO *postItemsDao.ItemDAO, purchaseDAO *purchaseItemDao.PurchaseDAO) *BlockchainUsecase {
	return &BlockchainUsecase{
		itemDAO:     itemDAO,
		purchaseDAO: purchaseDAO,
	}
}

// HandleItemListed はonchainで商品が登録された際に呼ばれる
// onchainのイベントから商品情報を取得してDBに挿入する
func (uc *BlockchainUsecase) HandleItemListed(chainItemID int64, tokenID int64, title string, priceWei string, explanation string, imageURL string, uid string, category string, seller string, createdAt int64, txHash string) error {
	log.Printf("HandleItemListed called: chain_item_id=%d, title=%s, uid=%s, seller=%s, price_wei=%s", chainItemID, title, uid, seller, priceWei)
	
	// バリデーション
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if uid == "" {
		log.Printf("WARNING: uid is empty for chain_item_id=%d, title=%s, seller=%s", chainItemID, title, seller)
		// uidが空でも処理を続行（sellerアドレスから推測できる可能性がある）
		// ただし、データベースの制約でエラーになる可能性がある
	}
	if seller == "" {
		return fmt.Errorf("seller address is required")
	}

	// 既にchain_item_idで商品が存在するか確認
	existingItemID, err := uc.itemDAO.FindItemByChainItemID(chainItemID)
	if err == nil && existingItemID > 0 {
		log.Printf("Item with chain_item_id=%d already exists (item_id=%d), updating...", chainItemID, existingItemID)
		// 既に存在する場合は更新のみ
		if err := uc.itemDAO.UpdateChainItemID(existingItemID, chainItemID, seller, tokenID); err != nil {
			return fmt.Errorf("failed to update chain_item_id: %w", err)
		}
		log.Printf("Successfully updated existing item (item_id=%d)", existingItemID)
		return nil
	}
	log.Printf("No existing item found with chain_item_id=%d", chainItemID)

	// 価格をWeiから円に変換（1円 = 0.000001 ETH = 10^12 Wei）
	priceInt := 0
	if priceWei != "" {
		priceBig, ok := new(big.Int).SetString(priceWei, 10)
		if ok {
			// WeiをETHに変換（1 ETH = 10^18 Wei）
			// ETHを円に変換（1円 = 0.000001 ETH = 10^-6 ETH）
			// つまり、1円 = 10^12 Wei
			// priceWei / 10^12 = 円
			ethPerJpy := big.NewInt(1e12) // 1円 = 10^12 Wei
			priceInt = int(new(big.Int).Div(priceBig, ethPerJpy).Int64())
		}
	}

	// 新規商品を作成（chain_item_idを含む）
	imageURLs := []string{}
	if imageURL != "" {
		imageURLs = append(imageURLs, imageURL)
	}

	// 既存の商品（uidとtitleでマッチング、chain_item_idがNULL）を検索
	existingItemID, err = uc.itemDAO.FindItemByUidAndTitle(uid, title)
	if err == nil && existingItemID > 0 {
		log.Printf("Found existing item by uid and title (item_id=%d), linking chain_item_id...", existingItemID)
		// 既存の商品にchain_item_idを関連付ける
		if err := uc.itemDAO.UpdateChainItemID(existingItemID, chainItemID, seller, tokenID); err != nil {
			return fmt.Errorf("failed to update chain_item_id: %w", err)
		}
		log.Printf("Successfully linked chain_item_id=%d to existing item (item_id=%d)", chainItemID, existingItemID)
		return nil
	}
	log.Printf("No existing item found by uid=%s and title=%s, creating new item...", uid, title)

	// InsertItemWithChainIDを使用してchain_item_idを含めて挿入
	log.Printf("Inserting new item: title=%s, price=%d, chain_item_id=%d, uid=%s", title, priceInt, chainItemID, uid)
	if err := uc.itemDAO.InsertItemWithChainID(title, priceInt, explanation, imageURLs, uid, "listed", category, chainItemID, seller, tokenID); err != nil {
		log.Printf("Error inserting item: %v", err)
		return fmt.Errorf("failed to create item: %w", err)
	}
	log.Printf("Successfully created new item with chain_item_id=%d", chainItemID)

	return nil
}

// HandleItemPurchased はonchainで商品が購入された際に呼ばれる
func (uc *BlockchainUsecase) HandleItemPurchased(chainItemID int64, buyer string, priceWei string, tokenID int64, txHash string) error {
	log.Printf("HandleItemPurchased called: chain_item_id=%d, buyer=%s, txHash=%s", chainItemID, buyer, txHash)

	// chain_item_idで商品を検索
	itemID, err := uc.itemDAO.FindItemByChainItemID(chainItemID)
	if err != nil {
		return fmt.Errorf("item not found for chain_item_id %d: %w", chainItemID, err)
	}
	log.Printf("Found item_id=%d for chain_item_id=%d", itemID, chainItemID)

	// buyerアドレスからUIDを取得（見つからない場合は空文字列）
	buyerUID := ""
	uid, err := uc.purchaseDAO.GetUIDByWalletAddress(buyer)
	if err == nil {
		buyerUID = uid
		log.Printf("Found buyer UID=%s for address=%s", buyerUID, buyer)
	} else {
		log.Printf("Buyer UID not found for address=%s (this is OK if user not registered)", buyer)
	}

	// 購入状態を更新（buyer_addressも保存）
	if err := uc.purchaseDAO.UpdatePurchaseStatus(int(itemID), buyerUID, buyer); err != nil {
		return fmt.Errorf("failed to update purchase status: %w", err)
	}

	log.Printf("Successfully updated purchase status: item_id=%d, chain_item_id=%d, buyer=%s", itemID, chainItemID, buyer)
	return nil
}

// HandleReceiptConfirmed はonchainで商品受け取り確認された際に呼ばれる
func (uc *BlockchainUsecase) HandleReceiptConfirmed(chainItemID int64, buyer string, seller string, priceWei string, txHash string) error {
	log.Printf("HandleReceiptConfirmed called: chain_item_id=%d, buyer=%s, seller=%s, txHash=%s", chainItemID, buyer, seller, txHash)

	// ステータスをcompletedに更新
	if err := uc.purchaseDAO.UpdateToCompleted(chainItemID); err != nil {
		return fmt.Errorf("failed to update status to completed: %w", err)
	}

	log.Printf("Successfully updated status to completed: chain_item_id=%d", chainItemID)
	return nil
}

// HandleItemCancelled はonchainで商品がキャンセルされた際に呼ばれる
func (uc *BlockchainUsecase) HandleItemCancelled(chainItemID int64, seller string, txHash string) error {
	log.Printf("HandleItemCancelled called: chain_item_id=%d, seller=%s, txHash=%s", chainItemID, seller, txHash)

	// ステータスをcancelledに更新
	if err := uc.purchaseDAO.UpdateToCancelled(chainItemID); err != nil {
		return fmt.Errorf("failed to update status to cancelled: %w", err)
	}

	log.Printf("Successfully updated status to cancelled: chain_item_id=%d", chainItemID)
	return nil
}

