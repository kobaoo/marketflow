package web

import (
	"encoding/json"
	"net/http"
)

// --------- /prices/latest ---------

func (h *Handler) getLatestPrice(w http.ResponseWriter, r *http.Request) {
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

// --------- /prices/highest ---------

func (h *Handler) getHighestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	if period, ok, err := parsePeriod(r); err != nil {
		http.Error(w, "Invalid period format", http.StatusBadRequest)
		return
	} else if ok {
		price := h.marketDataService.GetHighestPriceBySymbolAndPeriod(r.Context(), symbol, period)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(price)
		return
	}

	price := h.marketDataService.GetHighestPriceBySymbol(r.Context(), symbol)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(price)
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

	if period, ok, err := parsePeriod(r); err != nil {
		http.Error(w, "Invalid period format", http.StatusBadRequest)
		return
	} else if ok {
		price := h.marketDataService.GetHighestPriceBySymbolAndPeriodAndExchange(r.Context(), symbol, exchange, period)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(price)
		return
	}

	price := h.marketDataService.GetHighestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(price)
}

// --------- /prices/lowest ---------

func (h *Handler) getLowestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	if period, ok, err := parsePeriod(r); err != nil {
		http.Error(w, "Invalid period format", http.StatusBadRequest)
		return
	} else if ok {
		price := h.marketDataService.GetLowestPriceBySymbolAndPeriod(r.Context(), symbol, period)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(price)
		return
	}

	price := h.marketDataService.GetLowestPriceBySymbol(r.Context(), symbol)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(price)
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

	if period, ok, err := parsePeriod(r); err != nil {
		http.Error(w, "Invalid period format", http.StatusBadRequest)
		return
	} else if ok {
		price := h.marketDataService.GetLowestPriceBySymbolAndPeriodAndExchange(r.Context(), symbol, exchange, period)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(price)
		return
	}

	price := h.marketDataService.GetLowestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(price)
}

// --------- /prices/average ---------

func (h *Handler) getAveragePrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	if symbol == "" {
		http.Error(w, "Symbol parameter is required", http.StatusBadRequest)
		return
	}

	if period, ok, err := parsePeriod(r); err != nil {
		http.Error(w, "Invalid period format", http.StatusBadRequest)
		return
	} else if ok {
		price := h.marketDataService.GetAvgPriceBySymbolAndPeriod(r.Context(), symbol, period)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(price)
		return
	}

	price := h.marketDataService.GetAvgPriceBySymbol(r.Context(), symbol)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(price)
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

	if period, ok, err := parsePeriod(r); err != nil {
		http.Error(w, "Invalid period format", http.StatusBadRequest)
		return
	} else if ok {
		price := h.marketDataService.GetAvgPriceBySymbolAndPeriodAndExchange(r.Context(), symbol, exchange, period)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(price)
		return
	}

	price := h.marketDataService.GetAvgPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(price)
}
