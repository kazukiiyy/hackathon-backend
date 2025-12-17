package purchaseItem

import (
	dao "uttc-hackathon-backend/dao/purchaseItem"
)

type PurchaseUsecase struct {
	purchaseDAO dao.PurchaseDAOInterface
}

func NewPurchaseUsecase(purchaseDAO dao.PurchaseDAOInterface) *PurchaseUsecase {
	return &PurchaseUsecase{purchaseDAO: purchaseDAO}
}

func (u *PurchaseUsecase) PurchaseItem(itemID int, buyerUID string) error {
	// 従来の購入フロー（cash購入）ではbuyer_addressは空文字列
	return u.purchaseDAO.UpdatePurchaseStatus(itemID, buyerUID, "")
}

func (u *PurchaseUsecase) GetPurchasedItems(buyerUID string) ([]*dao.PurchasedItem, error) {
	return u.purchaseDAO.GetPurchasedItems(buyerUID)
}
