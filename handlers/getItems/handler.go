package getItems

import (
	"encoding/json"
	"net/http"
	"strconv"
	"uttc-hackathon-backend/usecase/getItems"
)

type ItemHandler struct {
	getItemUc *getItems.ItemUsecase
}

func NewItemHandler(u *getItems.ItemUsecase) *ItemHandler {
	return &ItemHandler{getItemUc: u}
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		writeJSONError(w, "category is required", http.StatusBadRequest)
		return
	}

	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	items, err := h.getItemUc.GetItemsByCategory(category, page, limit)
	if err != nil {
		writeJSONError(w, "Items not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		writeJSONError(w, "JSON encode error", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
