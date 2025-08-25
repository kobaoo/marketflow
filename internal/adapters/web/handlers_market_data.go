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
func (h *Handler) getLatestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}
	exchange := r.PathValue("exchange")
	if exchange == "" {
		http.Error(w, "Exchange parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetLatestPriceBySymbolAndExchange(r.Context(), symbol, exchange)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getHighestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}
	price := h.marketDataService.GetHighestPriceBySymbol(r.Context(), symbol)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getHighestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}
	exchange := r.PathValue("exchange")
	if exchange == "" {
		http.Error(w, "Exchange parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetHighestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getLowestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetLowestPriceBySymbol(r.Context(), symbol)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getLowestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}
	exchange := r.PathValue("exchange")
	if exchange == "" {
		http.Error(w, "Exchange parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetLowestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAveragePrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetAvgPriceBySymbol(r.Context(), symbol)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAveragePriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}
	exchange := r.PathValue("exchange")
	if exchange == "" {
		http.Error(w, "Exchange parameter is required", http.StatusBadRequest)
		return
	}

	price := h.marketDataService.GetAvgPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
