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

func (u *ItemUsecase) GetItemByID(id int) (*getItemDao.Item, error) {
	item, err := u.getItemDao.GetItemByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	return item, nil
}

func (u *ItemUsecase) GetItemsByUid(uid string) ([]*getItemDao.Item, error) {
	items, err := u.getItemDao.GetItemsByUid(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get items by uid: %w", err)
	}
	return items, nil
}

func (u *ItemUsecase) GetLatestItems(limit int) ([]*getItemDao.Item, error) {
	items, err := u.getItemDao.GetLatestItems(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest items: %w", err)
	}
	return items, nil
}
