package purchaseItem

import (
	"database/sql"
	"encoding/json"
	"log"
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
	path := r.URL.Path
	if !strings.HasSuffix(path, "/purchase") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid URL format. Expected: /items/{id}/purchase"})
		return
	}

	// /items/ を削除し、/purchase を削除
	path = strings.TrimPrefix(path, "/items/")
	path = strings.TrimSuffix(path, "/purchase")
	if path == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Item ID is required"})
		return
	}

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
		errMsg := err.Error()
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Item not found or already purchased"})
		} else if strings.Contains(errMsg, "seller cannot purchase their own item") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Seller cannot purchase their own item"})
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

// GET /purchases?buyer_uid=xxx&buyer_address=xxx
func (h *PurchaseHandler) GetPurchasedItems(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GetPurchasedItems] Request received: Method=%s, URL=%s", r.Method, r.URL.String())
	
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	buyerUID := r.URL.Query().Get("buyer_uid")
	buyerAddress := r.URL.Query().Get("buyer_address")
	log.Printf("[GetPurchasedItems] buyer_uid from query: %s, buyer_address: %s", buyerUID, buyerAddress)
	
	if buyerUID == "" && buyerAddress == "" {
		log.Printf("[GetPurchasedItems] Error: both buyer_uid and buyer_address are empty")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "buyer_uid or buyer_address is required"})
		return
	}

	items, err := h.usecase.GetPurchasedItems(buyerUID, buyerAddress)
	if err != nil {
		log.Printf("[GetPurchasedItems] Error from usecase: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get purchased items"})
		return
	}

	if items == nil {
		log.Printf("[GetPurchasedItems] Items is nil, setting to empty array")
		items = []*dao.PurchasedItem{}
	}

	log.Printf("[GetPurchasedItems] Returning %d items for buyer_uid=%s, buyer_address=%s", len(items), buyerUID, buyerAddress)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("[GetPurchasedItems] Error encoding response: %v", err)
	}
}
