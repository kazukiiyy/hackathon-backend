package postUser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uttc-hackathon-backend/usecase/postUser"
)

type UserHandler struct {
	postUserUc *postUser.UserUsecase
}

func NewUserHandler(u *postUser.UserUsecase) *UserHandler {
	return &UserHandler{postUserUc: u}
}

type RegisterRequest struct {
	Uid       string `json:"uid"`
	Nickname  string `json:"nickname"`
	Sex       string `json:"sex"`
	Birthyear int    `json:"birthyear"`
	Birthdate int    `json:"birthdate"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Uid == "" {
		writeJSONError(w, "uid is required", http.StatusBadRequest)
		return
	}

	if req.Nickname == "" {
		writeJSONError(w, "nickname is required", http.StatusBadRequest)
		return
	}

	response, err := h.postUserUc.RegisterUser(req.Uid, req.Nickname, req.Sex, req.Birthyear, req.Birthdate)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("ユーザー登録完了: %s\n", req.Nickname)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeJSONError(w, "JSON encode error", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
