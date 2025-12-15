package purchaseItem

import (
	dao "uttc-hackathon-backend/dao/purchaseItem"
)

type PurchaseUsecase struct {
	purchaseDAO *dao.PurchaseDAO
}

func NewPurchaseUsecase(purchaseDAO *dao.PurchaseDAO) *PurchaseUsecase {
	return &PurchaseUsecase{purchaseDAO: purchaseDAO}
}

func (u *PurchaseUsecase) PurchaseItem(itemID int, buyerUID string) error {
	return u.purchaseDAO.UpdatePurchaseStatus(itemID, buyerUID)
}

func (u *PurchaseUsecase) GetPurchasedItems(buyerUID string) ([]*dao.PurchasedItem, error) {
	return u.purchaseDAO.GetPurchasedItems(buyerUID)
}
