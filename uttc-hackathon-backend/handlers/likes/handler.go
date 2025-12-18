package likes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"uttc-hackathon-backend/usecase/likes"
)

type LikeHandler struct {
	likeUc *likes.LikeUsecase
}

func NewLikeHandler(u *likes.LikeUsecase) *LikeHandler {
	return &LikeHandler{likeUc: u}
}

type LikeRequest struct {
	ItemID int    `json:"item_id"`
	UID    string `json:"uid"`
}

type LikeStatusResponse struct {
	Liked bool `json:"liked"`
	Count int  `json:"count"`
}

// POST /likes - いいね追加
// DELETE /likes - いいね削除
func (h *LikeHandler) HandleLike(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.addLike(w, r)
	case http.MethodDelete:
		h.removeLike(w, r)
	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LikeHandler) addLike(w http.ResponseWriter, r *http.Request) {
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ItemID == 0 || req.UID == "" {
		writeJSONError(w, "item_id and uid are required", http.StatusBadRequest)
		return
	}

	err := h.likeUc.AddLike(req.ItemID, req.UID)
	if err != nil {
		// 重複エラーの場合は成功として扱う
		if strings.Contains(err.Error(), "Duplicate") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Already liked"})
			return
		}
		writeJSONError(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Liked successfully"})
}

func (h *LikeHandler) removeLike(w http.ResponseWriter, r *http.Request) {
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ItemID == 0 || req.UID == "" {
		writeJSONError(w, "item_id and uid are required", http.StatusBadRequest)
		return
	}

	err := h.likeUc.RemoveLike(req.ItemID, req.UID)
	if err != nil {
		writeJSONError(w, "Failed to remove like", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Unliked successfully"})
}

// GET /likes/status?item_id=123&uid=xxx - いいね状態とカウント取得
func (h *LikeHandler) GetLikeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemIDStr := r.URL.Query().Get("item_id")
	uid := r.URL.Query().Get("uid")

	if itemIDStr == "" {
		writeJSONError(w, "item_id is required", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		writeJSONError(w, "Invalid item_id", http.StatusBadRequest)
		return
	}

	count, err := h.likeUc.GetLikeCount(itemID)
	if err != nil {
		writeJSONError(w, "Failed to get like count", http.StatusInternalServerError)
		return
	}

	var liked bool
	if uid != "" {
		liked, err = h.likeUc.IsLiked(itemID, uid)
		if err != nil {
			writeJSONError(w, "Failed to check like status", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LikeStatusResponse{
		Liked: liked,
		Count: count,
	})
}

// GET /likes/user?uid=xxx - ユーザーがいいねした商品一覧
func (h *LikeHandler) GetUserLikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		writeJSONError(w, "uid is required", http.StatusBadRequest)
		return
	}

	itemIDs, err := h.likeUc.GetLikedItemsByUser(uid)
	if err != nil {
		writeJSONError(w, "Failed to get liked items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]int{"item_ids": itemIDs})
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
