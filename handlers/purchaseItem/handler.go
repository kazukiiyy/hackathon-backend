package purchaseItem

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	dao "uttc-hackathon-backend/dao/purchaseItem"
	uc "uttc-hackathon-backend/usecase/purchaseItem"
)

type PurchaseHandler struct {
	usecase *uc.PurchaseUsecase
}

func NewPurchaseHandler(usecase *uc.PurchaseUsecase) *PurchaseHandler {
	return &PurchaseHandler{usecase: usecase}
}

type PurchaseRequest struct {
	BuyerUID string `json:"buyer_uid"`
}

type PurchaseResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *PurchaseHandler) PurchaseItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	// URLからitem IDを取得 (/items/{id}/purchase)
	path := strings.TrimPrefix(r.URL.Path, "/items/")
	path = strings.TrimSuffix(path, "/purchase")
	itemID, err := strconv.Atoi(path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid item ID"})
		return
	}

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.BuyerUID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "buyer_uid is required"})
		return
	}

	err = h.usecase.PurchaseItem(itemID, req.BuyerUID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Item not found or already purchased"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to purchase item"})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PurchaseResponse{Message: "Purchase successful"})
}

// GET /purchases?buyer_uid=xxx
func (h *PurchaseHandler) GetPurchasedItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	buyerUID := r.URL.Query().Get("buyer_uid")
	if buyerUID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "buyer_uid is required"})
		return
	}

	items, err := h.usecase.GetPurchasedItems(buyerUID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get purchased items"})
		return
	}

	if items == nil {
		items = []*dao.PurchasedItem{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
