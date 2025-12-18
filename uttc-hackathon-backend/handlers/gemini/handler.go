package gemini

import (
	"encoding/json"
	"net/http"
	"uttc-hackathon-backend/usecase/gemini"
)

type GeminiHandler struct {
	geminiUC *gemini.GeminiUsecase
}

func NewGeminiHandler(uc *gemini.GeminiUsecase) *GeminiHandler {
	return &GeminiHandler{geminiUC: uc}
}

type GenerateContentRequest struct {
	Prompt   string `json:"prompt"`
	Protocol string `json:"protocol"` // 後で送るので今は空白
}

type GenerateContentResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func (h *GeminiHandler) GenerateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Prompt == "" {
		writeJSONError(w, "prompt is required", http.StatusBadRequest)
		return
	}

	result, err := h.geminiUC.GenerateContent(req.Prompt, req.Protocol)
	if err != nil {
		errorMsg := "Failed to generate content"
		if result != nil && result.Error != "" {
			errorMsg = result.Error
		} else {
			errorMsg = err.Error()
		}
		writeJSONError(w, errorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		writeJSONError(w, "JSON encode error", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
