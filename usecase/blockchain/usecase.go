package blockchain

import (
	"fmt"
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
	// 既にchain_item_idで商品が存在するか確認
	existingItemID, err := uc.itemDAO.FindItemByChainItemID(chainItemID)
	if err == nil && existingItemID > 0 {
		// 既に存在する場合は更新のみ
		if err := uc.itemDAO.UpdateChainItemID(existingItemID, chainItemID, seller, tokenID); err != nil {
			return fmt.Errorf("failed to update chain_item_id: %w", err)
		}
		return nil
	}

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
	existingItemID, err := uc.itemDAO.FindItemByUidAndTitle(uid, title)
	if err == nil && existingItemID > 0 {
		// 既存の商品にchain_item_idを関連付ける
		if err := uc.itemDAO.UpdateChainItemID(existingItemID, chainItemID, seller, tokenID); err != nil {
			return fmt.Errorf("failed to update chain_item_id: %w", err)
		}
		return nil
	}

	// InsertItemWithChainIDを使用してchain_item_idを含めて挿入
	if err := uc.itemDAO.InsertItemWithChainID(title, priceInt, explanation, imageURLs, uid, "listed", category, chainItemID, seller, tokenID); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

// HandleItemPurchased はonchainで商品が購入された際に呼ばれる
func (uc *BlockchainUsecase) HandleItemPurchased(chainItemID int64, buyer string, priceWei string, tokenID int64, txHash string) error {
	// chain_item_idで商品を検索
	itemID, err := uc.itemDAO.FindItemByChainItemID(chainItemID)
	if err != nil {
		return fmt.Errorf("item not found for chain_item_id %d: %w", chainItemID, err)
	}

	// 購入状態を更新（buyer_uidは空文字列として扱う、またはbuyerアドレスをUIDとして扱う）
	// 注意: buyerはアドレスなので、UIDに変換する必要がある場合がある
	// ここでは簡易的にbuyerアドレスをUIDとして扱う
	if err := uc.purchaseDAO.UpdatePurchaseStatus(int(itemID), buyer); err != nil {
		return fmt.Errorf("failed to update purchase status: %w", err)
	}

	return nil
}

