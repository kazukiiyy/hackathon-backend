package likes

import (
	"errors"
	"fmt"
	"testing"
)

// MockLikeDAO はテスト用のモックDAO
type MockLikeDAO struct {
	likes            map[string]bool // "itemID:uid" -> liked
	likeCounts       map[int]int     // itemID -> count
	addLikeErr       error
	removeLikeErr    error
	isLikedErr       error
	getLikeCountErr  error
	getLikedItemsErr error
}

func NewMockLikeDAO() *MockLikeDAO {
	return &MockLikeDAO{
		likes:      make(map[string]bool),
		likeCounts: make(map[int]int),
	}
}

func (m *MockLikeDAO) makeKey(itemID int, uid string) string {
	return fmt.Sprintf("%d:%s", itemID, uid)
}

func (m *MockLikeDAO) AddLike(itemID int, uid string) error {
	if m.addLikeErr != nil {
		return m.addLikeErr
	}
	key := m.makeKey(itemID, uid)
	if m.likes[key] {
		return errors.New("duplicate like")
	}
	m.likes[key] = true
	m.likeCounts[itemID]++
	return nil
}

func (m *MockLikeDAO) RemoveLike(itemID int, uid string) error {
	if m.removeLikeErr != nil {
		return m.removeLikeErr
	}
	key := m.makeKey(itemID, uid)
	if m.likes[key] {
		delete(m.likes, key)
		m.likeCounts[itemID]--
	}
	return nil
}

func (m *MockLikeDAO) IsLiked(itemID int, uid string) (bool, error) {
	if m.isLikedErr != nil {
		return false, m.isLikedErr
	}
	return m.likes[m.makeKey(itemID, uid)], nil
}

func (m *MockLikeDAO) GetLikeCount(itemID int) (int, error) {
	if m.getLikeCountErr != nil {
		return 0, m.getLikeCountErr
	}
	return m.likeCounts[itemID], nil
}

func (m *MockLikeDAO) GetLikedItemsByUser(uid string) ([]int, error) {
	if m.getLikedItemsErr != nil {
		return nil, m.getLikedItemsErr
	}
	var items []int
	for key, liked := range m.likes {
		if liked {
			var itemID int
			var keyUID string
			fmt.Sscanf(key, "%d:%s", &itemID, &keyUID)
			if keyUID == uid {
				items = append(items, itemID)
			}
		}
	}
	return items, nil
}

// TestAddLike_Success いいね追加成功
func TestAddLike_Success(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	err := usecase.AddLike(1, "user123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// いいねが追加されているか確認
	liked, _ := mockDAO.IsLiked(1, "user123")
	if !liked {
		t.Error("expected item to be liked")
	}

	// カウントが増えているか確認
	count, _ := mockDAO.GetLikeCount(1)
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

// TestAddLike_Duplicate 重複いいねエラー
func TestAddLike_Duplicate(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	// 1回目は成功
	_ = usecase.AddLike(1, "user123")

	// 2回目は失敗
	err := usecase.AddLike(1, "user123")
	if err == nil {
		t.Error("expected error for duplicate like")
	}
}

// TestAddLike_DAOError DAOエラー時
func TestAddLike_DAOError(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	mockDAO.addLikeErr = errors.New("database error")
	usecase := NewLikeUsecase(mockDAO)

	err := usecase.AddLike(1, "user123")
	if err == nil {
		t.Error("expected error")
	}
}

// TestRemoveLike_Success いいね削除成功
func TestRemoveLike_Success(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	// まずいいねを追加
	_ = usecase.AddLike(1, "user123")

	// いいねを削除
	err := usecase.RemoveLike(1, "user123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// いいねが削除されているか確認
	liked, _ := mockDAO.IsLiked(1, "user123")
	if liked {
		t.Error("expected item to not be liked")
	}

	// カウントが減っているか確認
	count, _ := mockDAO.GetLikeCount(1)
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

// TestRemoveLike_NotLiked いいねしていない商品の削除
func TestRemoveLike_NotLiked(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	// いいねしていない商品を削除してもエラーにならない
	err := usecase.RemoveLike(1, "user123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestIsLiked_True いいね済み
func TestIsLiked_True(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	_ = usecase.AddLike(1, "user123")

	liked, err := usecase.IsLiked(1, "user123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !liked {
		t.Error("expected liked to be true")
	}
}

// TestIsLiked_False いいねしていない
func TestIsLiked_False(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	liked, err := usecase.IsLiked(1, "user123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if liked {
		t.Error("expected liked to be false")
	}
}

// TestGetLikeCount カウント取得
func TestGetLikeCount(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	// 複数ユーザーがいいね
	_ = usecase.AddLike(1, "user1")
	_ = usecase.AddLike(1, "user2")
	_ = usecase.AddLike(1, "user3")

	count, err := usecase.GetLikeCount(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

// TestGetLikedItemsByUser ユーザーのいいね一覧
func TestGetLikedItemsByUser(t *testing.T) {
	mockDAO := NewMockLikeDAO()
	usecase := NewLikeUsecase(mockDAO)

	// ユーザーが複数商品にいいね
	_ = usecase.AddLike(1, "user123")
	_ = usecase.AddLike(2, "user123")
	_ = usecase.AddLike(3, "user123")
	// 別ユーザーのいいね
	_ = usecase.AddLike(4, "other_user")

	items, err := usecase.GetLikedItemsByUser("user123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}
