package likes

import (
	"fmt"
	likesDao "uttc-hackathon-backend/dao/likes"
)

type LikeUsecase struct {
	likeDao *likesDao.LikeDAO
}

func NewLikeUsecase(dao *likesDao.LikeDAO) *LikeUsecase {
	return &LikeUsecase{likeDao: dao}
}

func (u *LikeUsecase) AddLike(itemID int, uid string) error {
	err := u.likeDao.AddLike(itemID, uid)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	return nil
}

func (u *LikeUsecase) RemoveLike(itemID int, uid string) error {
	err := u.likeDao.RemoveLike(itemID, uid)
	if err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}
	return nil
}

func (u *LikeUsecase) IsLiked(itemID int, uid string) (bool, error) {
	liked, err := u.likeDao.IsLiked(itemID, uid)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %w", err)
	}
	return liked, nil
}

func (u *LikeUsecase) GetLikeCount(itemID int) (int, error) {
	count, err := u.likeDao.GetLikeCount(itemID)
	if err != nil {
		return 0, fmt.Errorf("failed to get like count: %w", err)
	}
	return count, nil
}

func (u *LikeUsecase) GetLikedItemsByUser(uid string) ([]int, error) {
	itemIDs, err := u.likeDao.GetLikedItemsByUser(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked items: %w", err)
	}
	return itemIDs, nil
}
