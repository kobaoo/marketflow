package web

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) getLatestPrice(w http.ResponseWriter, r *http.Request) {
	// Extract symbol from URL path
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetLatestPriceBySymbol(r.Context(), symbol)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getLatestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func getHighestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func getHighestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func getLowestPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func getLowestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func getAveragePrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}

func getAveragePriceFromExchange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "not implemented"})
}
