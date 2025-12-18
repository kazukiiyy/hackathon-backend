package purchaseItem

import (
	"errors"
	"testing"
	"time"

	dao "uttc-hackathon-backend/dao/purchaseItem"
)

// MockPurchaseDAO はテスト用のモックDAO
type MockPurchaseDAO struct {
	purchasedItems map[int]bool           // itemID -> purchased
	purchases      map[string][]*dao.PurchasedItem // buyerUID -> items
	updateErr      error
	getItemsErr    error
}

func NewMockPurchaseDAO() *MockPurchaseDAO {
	return &MockPurchaseDAO{
		purchasedItems: make(map[int]bool),
		purchases:      make(map[string][]*dao.PurchasedItem),
	}
}

func (m *MockPurchaseDAO) UpdatePurchaseStatus(itemID int, buyerUID string, buyerAddress string) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if m.purchasedItems[itemID] {
		return errors.New("item already purchased")
	}
	m.purchasedItems[itemID] = true

	// 購入履歴に追加
	item := &dao.PurchasedItem{
		ID:          itemID,
		Title:       "Test Item",
		Price:       1000,
		PurchasedAt: time.Now(),
	}
	m.purchases[buyerUID] = append(m.purchases[buyerUID], item)
	return nil
}

func (m *MockPurchaseDAO) GetUIDByWalletAddress(walletAddress string) (string, error) {
	// テスト用の簡易実装（実際の実装ではDBから取得）
	return "", errors.New("not found")
}

func (m *MockPurchaseDAO) GetPurchasedItems(buyerUID string, buyerAddress string) ([]*dao.PurchasedItem, error) {
	if m.getItemsErr != nil {
		return nil, m.getItemsErr
	}
	return m.purchases[buyerUID], nil
}

func (m *MockPurchaseDAO) GetWalletAddressByUID(uid string) (string, error) {
	// テスト用の簡易実装（実際の実装ではDBから取得）
	return "", errors.New("not found")
}

// TestPurchaseItem_Success 購入成功
func TestPurchaseItem_Success(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	usecase := NewPurchaseUsecase(mockDAO)

	err := usecase.PurchaseItem(1, "buyer123")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 購入済みになっているか確認
	if !mockDAO.purchasedItems[1] {
		t.Error("expected item to be purchased")
	}
}

// TestPurchaseItem_AlreadyPurchased 購入済み商品の再購入
func TestPurchaseItem_AlreadyPurchased(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	usecase := NewPurchaseUsecase(mockDAO)

	// 1回目は成功
	_ = usecase.PurchaseItem(1, "buyer123")

	// 2回目は失敗（別のユーザーでも）
	err := usecase.PurchaseItem(1, "buyer456")
	if err == nil {
		t.Error("expected error for already purchased item")
	}
}

// TestPurchaseItem_DAOError DAOエラー時
func TestPurchaseItem_DAOError(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	mockDAO.updateErr = errors.New("database error")
	usecase := NewPurchaseUsecase(mockDAO)

	err := usecase.PurchaseItem(1, "buyer123")
	if err == nil {
		t.Error("expected error")
	}
}

// TestGetPurchasedItems_Success 購入履歴取得成功
func TestGetPurchasedItems_Success(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	usecase := NewPurchaseUsecase(mockDAO)

	// 複数商品を購入
	_ = usecase.PurchaseItem(1, "buyer123")
	_ = usecase.PurchaseItem(2, "buyer123")
	_ = usecase.PurchaseItem(3, "buyer123")

	items, err := usecase.GetPurchasedItems("buyer123", "")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

// TestGetPurchasedItems_Empty 購入履歴なし
func TestGetPurchasedItems_Empty(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	usecase := NewPurchaseUsecase(mockDAO)

	items, err := usecase.GetPurchasedItems("buyer123", "")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

// TestGetPurchasedItems_DAOError DAOエラー時
func TestGetPurchasedItems_DAOError(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	mockDAO.getItemsErr = errors.New("database error")
	usecase := NewPurchaseUsecase(mockDAO)

	_, err := usecase.GetPurchasedItems("buyer123")
	if err == nil {
		t.Error("expected error")
	}
}

// TestGetPurchasedItems_DifferentUsers 異なるユーザーの購入履歴は分離
func TestGetPurchasedItems_DifferentUsers(t *testing.T) {
	mockDAO := NewMockPurchaseDAO()
	usecase := NewPurchaseUsecase(mockDAO)

	// 異なるユーザーがそれぞれ購入
	_ = usecase.PurchaseItem(1, "buyer1")
	_ = usecase.PurchaseItem(2, "buyer1")
	_ = usecase.PurchaseItem(3, "buyer2")

	items1, _ := usecase.GetPurchasedItems("buyer1", "")
	items2, _ := usecase.GetPurchasedItems("buyer2", "")

	if len(items1) != 2 {
		t.Errorf("expected 2 items for buyer1, got %d", len(items1))
	}
	if len(items2) != 1 {
		t.Errorf("expected 1 item for buyer2, got %d", len(items2))
	}
}
