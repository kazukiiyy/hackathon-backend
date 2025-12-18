package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
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
	log.Printf("HandleItemListed called: method=%s, path=%s", r.Method, r.URL.Path)
	
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
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

	// リクエストボディを読み取る前にログ出力
	log.Printf("Reading request body...")
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received ItemListed event: chain_item_id=%d, title=%s, uid=%s, seller=%s, price_wei=%s", req.ChainItemID, req.Title, req.UID, req.Seller, req.PriceWei)
	
	// 空の値チェック
	if req.UID == "" {
		log.Printf("WARNING: uid is empty in request")
	}
	if req.Title == "" {
		log.Printf("WARNING: title is empty in request")
	}
	if req.Seller == "" {
		log.Printf("WARNING: seller is empty in request")
	}

	if err := h.blockchainUC.HandleItemListed(req.ChainItemID, req.TokenID, req.Title, req.PriceWei, req.Explanation, req.ImageURL, req.UID, req.Category, req.Seller, req.CreatedAt, req.TxHash); err != nil {
		log.Printf("Error processing ItemListed event: %v", err)
		writeJSONError(w, fmt.Sprintf("Failed to process item listed event: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully processed ItemListed event for chain_item_id=%d", req.ChainItemID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item listed event processed successfully"})
}

// HandleItemPurchased はonchainサービスからItemPurchasedイベントを受け取る
func (h *BlockchainHandler) HandleItemPurchased(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleItemPurchased called: method=%s, path=%s", r.Method, r.URL.Path)
	
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
		log.Printf("Error decoding request body: %v", err)
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received ItemPurchased event: chain_item_id=%d, buyer=%s, txHash=%s", req.ChainItemID, req.Buyer, req.TxHash)

	if err := h.blockchainUC.HandleItemPurchased(req.ChainItemID, req.Buyer, req.PriceWei, req.TokenID, req.TxHash); err != nil {
		log.Printf("Error processing ItemPurchased event: %v", err)
		writeJSONError(w, fmt.Sprintf("Failed to process item purchased event: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully processed ItemPurchased event for chain_item_id=%d", req.ChainItemID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item purchased event processed successfully"})
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

