package messages

import (
	"database/sql"
	"time"
)

type Message struct {
	ID          int       `json:"id"`
	SenderUID   string    `json:"sender_uid"`
	ReceiverUID string    `json:"receiver_uid"`
	Content     string    `json:"content"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type MessageDAO struct {
	db *sql.DB
}

func NewMessageDAO(db *sql.DB) *MessageDAO {
	return &MessageDAO{db: db}
}

// 相手とのメッセージ一覧を取得
func (d *MessageDAO) GetMessagesByPartner(myUID, partnerUID string) ([]*Message, error) {
	query := `
		SELECT id, sender_uid, receiver_uid, content, is_read, created_at
		FROM messages
		WHERE (sender_uid = ? AND receiver_uid = ?)
		   OR (sender_uid = ? AND receiver_uid = ?)
		ORDER BY created_at ASC
	`
	rows, err := d.db.Query(query, myUID, partnerUID, partnerUID, myUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.SenderUID, &msg.ReceiverUID, &msg.Content, &msg.IsRead, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, rows.Err()
}

// メッセージ送信
func (d *MessageDAO) CreateMessage(senderUID, receiverUID, content string) (*Message, error) {
	query := "INSERT INTO messages (sender_uid, receiver_uid, content) VALUES (?, ?, ?)"
	result, err := d.db.Exec(query, senderUID, receiverUID, content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:          int(id),
		SenderUID:   senderUID,
		ReceiverUID: receiverUID,
		Content:     content,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}, nil
}

// 相手からのメッセージを既読にする
func (d *MessageDAO) MarkAsRead(myUID, partnerUID string) error {
	query := "UPDATE messages SET is_read = true WHERE sender_uid = ? AND receiver_uid = ? AND is_read = false"
	_, err := d.db.Exec(query, partnerUID, myUID)
	return err
}

// Conversation はやり取りしている相手との会話情報
type Conversation struct {
	PartnerUID     string    `json:"partner_uid"`
	LastMessage    string    `json:"last_message"`
	LastMessageAt  time.Time `json:"last_message_at"`
	UnreadCount    int       `json:"unread_count"`
}

// やり取りしている相手一覧を取得
func (d *MessageDAO) GetConversations(myUID string) ([]*Conversation, error) {
	query := `
		SELECT
			partner_uid,
			last_message,
			last_message_at,
			unread_count
		FROM (
			SELECT
				CASE
					WHEN sender_uid = ? THEN receiver_uid
					ELSE sender_uid
				END as partner_uid,
				content as last_message,
				created_at as last_message_at,
				(SELECT COUNT(*) FROM messages m2
				 WHERE m2.sender_uid = CASE WHEN m1.sender_uid = ? THEN m1.receiver_uid ELSE m1.sender_uid END
				 AND m2.receiver_uid = ?
				 AND m2.is_read = false) as unread_count,
				ROW_NUMBER() OVER (
					PARTITION BY CASE WHEN sender_uid = ? THEN receiver_uid ELSE sender_uid END
					ORDER BY created_at DESC
				) as rn
			FROM messages m1
			WHERE sender_uid = ? OR receiver_uid = ?
		) sub
		WHERE rn = 1
		ORDER BY last_message_at DESC
	`
	rows, err := d.db.Query(query, myUID, myUID, myUID, myUID, myUID, myUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.PartnerUID, &conv.LastMessage, &conv.LastMessageAt, &conv.UnreadCount)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, &conv)
	}

	return conversations, rows.Err()
}
