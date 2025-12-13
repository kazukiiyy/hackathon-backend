package purchaseItem

type PurchaseDAO interface {
	UpdatePurchaseStatus(itemID int, buyerUID string) error
}

type PurchaseUsecase struct {
	dao PurchaseDAO
}

func NewPurchaseUsecase(dao PurchaseDAO) *PurchaseUsecase {
	return &PurchaseUsecase{dao: dao}
}

func (u *PurchaseUsecase) PurchaseItem(itemID int, buyerUID string) error {
	return u.dao.UpdatePurchaseStatus(itemID, buyerUID)
}
