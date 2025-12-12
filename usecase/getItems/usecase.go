package getItems

import (
	"fmt"
	getItemDao "uttc-hackathon-backend/dao/getItems"
)

type ItemUsecase struct {
	getItemDao *getItemDao.ItemDAO
}

func NewItemUsecase(dao *getItemDao.ItemDAO) *ItemUsecase {
	return &ItemUsecase{getItemDao: dao}
}

func (u *ItemUsecase) GetItemsByCategory(category string, page, limit int) ([]*getItemDao.Item, error) {
	items, err := u.getItemDao.GetItemsByCategory(category, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	return items, nil
}
