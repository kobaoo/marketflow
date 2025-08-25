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
func getLatestPriceFromExchange(w http.ResponseWriter, r *http.Request) {}

func getHighestPrice(w http.ResponseWriter, r *http.Request)                       {}
func getHighestPriceFromExchange(w http.ResponseWriter, r *http.Request)           {}
func getHighestPriceWithPeriod(w http.ResponseWriter, r *http.Request)             {}
func getHighestPriceFromExchangeWithPeriod(w http.ResponseWriter, r *http.Request) {}

func getLowestPrice(w http.ResponseWriter, r *http.Request)                       {}
func getLowestPriceFromExchange(w http.ResponseWriter, r *http.Request)           {}
func getLowestPriceWithPeriod(w http.ResponseWriter, r *http.Request)             {}
func getLowestPriceFromExchangeWithPeriod(w http.ResponseWriter, r *http.Request) {}

func getAveragePrice(w http.ResponseWriter, r *http.Request)                       {}
func getAveragePriceFromExchange(w http.ResponseWriter, r *http.Request)           {}
func getAveragePriceWithPeriod(w http.ResponseWriter, r *http.Request)             {}
func getAveragePriceFromExchangeWithPeriod(w http.ResponseWriter, r *http.Request) {}
