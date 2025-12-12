package getItems

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}

	items, err := h.getItemUc.GetItemsByCategory(category)
	if err != nil {
		http.Error(w, "Items not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
}
