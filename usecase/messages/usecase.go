package messages

import (
	dao "uttc-hackathon-backend/dao/messages"
)

type MessageUsecase struct {
	messageDAO *dao.MessageDAO
}

func NewMessageUsecase(messageDAO *dao.MessageDAO) *MessageUsecase {
	return &MessageUsecase{messageDAO: messageDAO}
}

func (u *MessageUsecase) GetMessages(myUID, partnerUID string) ([]*dao.Message, error) {
	return u.messageDAO.GetMessagesByPartner(myUID, partnerUID)
}

func (u *MessageUsecase) SendMessage(senderUID, receiverUID, content string) (*dao.Message, error) {
	return u.messageDAO.CreateMessage(senderUID, receiverUID, content)
}

func (u *MessageUsecase) MarkAsRead(myUID, partnerUID string) error {
	return u.messageDAO.MarkAsRead(myUID, partnerUID)
}
