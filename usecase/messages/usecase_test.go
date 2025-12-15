package messages

import (
	"errors"
	"testing"
	"time"

	dao "uttc-hackathon-backend/dao/messages"
)

// MockMessageDAO はテスト用のモックDAO
type MockMessageDAO struct {
	messages       []*dao.Message
	nextID         int
	createErr      error
	getMessagesErr error
	markAsReadErr  error
	getConvErr     error
}

func NewMockMessageDAO() *MockMessageDAO {
	return &MockMessageDAO{
		messages: make([]*dao.Message, 0),
		nextID:   1,
	}
}

func (m *MockMessageDAO) GetMessagesByPartner(myUID, partnerUID string) ([]*dao.Message, error) {
	if m.getMessagesErr != nil {
		return nil, m.getMessagesErr
	}

	var result []*dao.Message
	for _, msg := range m.messages {
		if (msg.SenderUID == myUID && msg.ReceiverUID == partnerUID) ||
			(msg.SenderUID == partnerUID && msg.ReceiverUID == myUID) {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (m *MockMessageDAO) CreateMessage(senderUID, receiverUID, content string) (*dao.Message, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}

	msg := &dao.Message{
		ID:          m.nextID,
		SenderUID:   senderUID,
		ReceiverUID: receiverUID,
		Content:     content,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}
	m.nextID++
	m.messages = append(m.messages, msg)
	return msg, nil
}

func (m *MockMessageDAO) MarkAsRead(myUID, partnerUID string) error {
	if m.markAsReadErr != nil {
		return m.markAsReadErr
	}

	for _, msg := range m.messages {
		if msg.SenderUID == partnerUID && msg.ReceiverUID == myUID && !msg.IsRead {
			msg.IsRead = true
		}
	}
	return nil
}

func (m *MockMessageDAO) GetConversations(myUID string) ([]*dao.Conversation, error) {
	if m.getConvErr != nil {
		return nil, m.getConvErr
	}

	// 簡易的に相手ごとの最新メッセージを返す
	partners := make(map[string]*dao.Conversation)
	for _, msg := range m.messages {
		var partnerUID string
		if msg.SenderUID == myUID {
			partnerUID = msg.ReceiverUID
		} else if msg.ReceiverUID == myUID {
			partnerUID = msg.SenderUID
		} else {
			continue
		}

		if _, ok := partners[partnerUID]; !ok {
			partners[partnerUID] = &dao.Conversation{
				PartnerUID:    partnerUID,
				LastMessage:   msg.Content,
				LastMessageAt: msg.CreatedAt,
				UnreadCount:   0,
			}
		} else {
			conv := partners[partnerUID]
			if msg.CreatedAt.After(conv.LastMessageAt) {
				conv.LastMessage = msg.Content
				conv.LastMessageAt = msg.CreatedAt
			}
		}

		// 未読カウント
		if msg.SenderUID == partnerUID && msg.ReceiverUID == myUID && !msg.IsRead {
			partners[partnerUID].UnreadCount++
		}
	}

	var result []*dao.Conversation
	for _, conv := range partners {
		result = append(result, conv)
	}
	return result, nil
}

// TestSendMessage_Success メッセージ送信成功
func TestSendMessage_Success(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	msg, err := usecase.SendMessage("sender123", "receiver456", "Hello!")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if msg == nil {
		t.Fatal("expected message, got nil")
	}
	if msg.Content != "Hello!" {
		t.Errorf("expected content 'Hello!', got '%s'", msg.Content)
	}
	if msg.SenderUID != "sender123" {
		t.Errorf("expected sender 'sender123', got '%s'", msg.SenderUID)
	}
	if msg.IsRead {
		t.Error("expected message to be unread")
	}
}

// TestSendMessage_DAOError DAOエラー時
func TestSendMessage_DAOError(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	mockDAO.createErr = errors.New("database error")
	usecase := NewMessageUsecase(mockDAO)

	_, err := usecase.SendMessage("sender123", "receiver456", "Hello!")
	if err == nil {
		t.Error("expected error")
	}
}

// TestGetMessages_Success メッセージ取得成功
func TestGetMessages_Success(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	// メッセージを送信
	_, _ = usecase.SendMessage("user1", "user2", "Hello")
	_, _ = usecase.SendMessage("user2", "user1", "Hi there")
	_, _ = usecase.SendMessage("user1", "user2", "How are you?")

	// メッセージを取得
	messages, err := usecase.GetMessages("user1", "user2")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
	}
}

// TestGetMessages_Empty メッセージなし
func TestGetMessages_Empty(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	messages, err := usecase.GetMessages("user1", "user2")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(messages))
	}
}

// TestGetMessages_OnlyBetweenPartners 指定した相手とのメッセージのみ取得
func TestGetMessages_OnlyBetweenPartners(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	// 異なる相手とのメッセージ
	_, _ = usecase.SendMessage("user1", "user2", "To user2")
	_, _ = usecase.SendMessage("user1", "user3", "To user3")
	_, _ = usecase.SendMessage("user2", "user1", "From user2")

	// user1とuser2の会話のみ取得
	messages, _ := usecase.GetMessages("user1", "user2")
	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}
}

// TestMarkAsRead_Success 既読更新成功
func TestMarkAsRead_Success(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	// 相手からのメッセージを作成
	_, _ = usecase.SendMessage("user2", "user1", "Hello")
	_, _ = usecase.SendMessage("user2", "user1", "Are you there?")

	// 既読にする
	err := usecase.MarkAsRead("user1", "user2")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// メッセージが既読になっているか確認
	messages, _ := usecase.GetMessages("user1", "user2")
	for _, msg := range messages {
		if !msg.IsRead {
			t.Error("expected all messages to be read")
		}
	}
}

// TestGetConversations_Success 会話一覧取得成功
func TestGetConversations_Success(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	// 複数の相手とやり取り
	_, _ = usecase.SendMessage("user1", "user2", "Hello user2")
	_, _ = usecase.SendMessage("user1", "user3", "Hello user3")
	_, _ = usecase.SendMessage("user4", "user1", "Hello from user4")

	conversations, err := usecase.GetConversations("user1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(conversations) != 3 {
		t.Errorf("expected 3 conversations, got %d", len(conversations))
	}
}

// TestGetConversations_UnreadCount 未読カウント
func TestGetConversations_UnreadCount(t *testing.T) {
	mockDAO := NewMockMessageDAO()
	usecase := NewMessageUsecase(mockDAO)

	// 相手からの未読メッセージ
	_, _ = usecase.SendMessage("user2", "user1", "Message 1")
	_, _ = usecase.SendMessage("user2", "user1", "Message 2")
	_, _ = usecase.SendMessage("user2", "user1", "Message 3")

	conversations, _ := usecase.GetConversations("user1")
	if len(conversations) != 1 {
		t.Fatalf("expected 1 conversation, got %d", len(conversations))
	}
	if conversations[0].UnreadCount != 3 {
		t.Errorf("expected 3 unread, got %d", conversations[0].UnreadCount)
	}
}
