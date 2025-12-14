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
