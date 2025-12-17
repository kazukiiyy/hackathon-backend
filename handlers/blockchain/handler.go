package blockchain

import (
	"encoding/json"
	"net/http"
	"uttc-hackathon-backend/usecase/blockchain"
)

type BlockchainHandler struct {
	blockchainUC *blockchain.BlockchainUsecase
}

func NewBlockchainHandler(uc *blockchain.BlockchainUsecase) *BlockchainHandler {
	return &BlockchainHandler{blockchainUC: uc}
}

// HandleItemListed はonchainサービスからItemListedイベントを受け取る
func (h *BlockchainHandler) HandleItemListed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChainItemID int64  `json:"chain_item_id"`
		TokenID     int64  `json:"token_id"`
		Title       string `json:"title"`
		PriceWei    string `json:"price_wei"`
		Explanation string `json:"explanation"`
		ImageURL    string `json:"image_url"`
		UID         string `json:"uid"`
		Category    string `json:"category"`
		Seller      string `json:"seller"`
		CreatedAt   int64  `json:"created_at"`
		TxHash      string `json:"tx_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.blockchainUC.HandleItemListed(req.ChainItemID, req.TokenID, req.Title, req.PriceWei, req.Explanation, req.ImageURL, req.UID, req.Category, req.Seller, req.CreatedAt, req.TxHash); err != nil {
		writeJSONError(w, "Failed to process item listed event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item listed event processed successfully"})
}

// HandleItemPurchased はonchainサービスからItemPurchasedイベントを受け取る
func (h *BlockchainHandler) HandleItemPurchased(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChainItemID int64  `json:"chain_item_id"`
		Buyer       string `json:"buyer"`
		PriceWei    string `json:"price_wei"`
		TokenID     int64  `json:"token_id"`
		TxHash      string `json:"tx_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.blockchainUC.HandleItemPurchased(req.ChainItemID, req.Buyer, req.PriceWei, req.TokenID, req.TxHash); err != nil {
		writeJSONError(w, "Failed to process item purchased event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item purchased event processed successfully"})
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

