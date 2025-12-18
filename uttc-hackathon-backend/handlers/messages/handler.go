package messages

import (
	"encoding/json"
	"net/http"
	dao "uttc-hackathon-backend/dao/messages"
	uc "uttc-hackathon-backend/usecase/messages"
)

type MessageHandler struct {
	usecase *uc.MessageUsecase
}

func NewMessageHandler(usecase *uc.MessageUsecase) *MessageHandler {
	return &MessageHandler{usecase: usecase}
}

type SendMessageRequest struct {
	SenderUID   string `json:"sender_uid"`
	ReceiverUID string `json:"receiver_uid"`
	Content     string `json:"content"`
}

type MarkReadRequest struct {
	MyUID      string `json:"my_uid"`
	PartnerUID string `json:"partner_uid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// GET /messages?my_uid=xxx&partner_uid=yyy
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	myUID := r.URL.Query().Get("my_uid")
	partnerUID := r.URL.Query().Get("partner_uid")

	if myUID == "" || partnerUID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "my_uid and partner_uid are required"})
		return
	}

	messages, err := h.usecase.GetMessages(myUID, partnerUID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get messages"})
		return
	}

	if messages == nil {
		messages = []*dao.Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// POST /messages
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.SenderUID == "" || req.ReceiverUID == "" || req.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "sender_uid, receiver_uid, and content are required"})
		return
	}

	message, err := h.usecase.SendMessage(req.SenderUID, req.ReceiverUID, req.Content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to send message"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// PUT /messages/read
func (h *MessageHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	var req MarkReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.MyUID == "" || req.PartnerUID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "my_uid and partner_uid are required"})
		return
	}

	err := h.usecase.MarkAsRead(req.MyUID, req.PartnerUID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to mark as read"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Marked as read"})
}

// GET /messages/conversations?uid=xxx
func (h *MessageHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "uid is required"})
		return
	}

	conversations, err := h.usecase.GetConversations(uid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get conversations"})
		return
	}

	if conversations == nil {
		conversations = []*dao.Conversation{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}
