package web

import (
	"log/slog"
	"marketflow/internal/domain"
	"net/http"
)

func (h *Handler) getLatestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	price, err := h.marketDataService.GetLatestPriceBySymbol(r.Context(), symbol)
	if err != nil {
		slog.Error("Failed to get latest price by symbol", "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("Successfully got latest price by symbol", "symbol", symbol)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getLatestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.PathValue("exchange")

	price, err := h.marketDataService.GetLatestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	if err != nil {
		slog.Error("Failed to get latest price by symbol and exchange", "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("Successfully got latest price by symbol and exchange", "symbol", symbol, "exchange", exchange)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getHighestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	period, ok, err := parsePeriod(r)
	if err != nil {
		slog.Error("Invalid period format", "symbol", symbol, "err", err)
		ResponseError(w, http.StatusBadRequest, "Invalid period format")
		return
	}

	if ok {
		price, err := h.marketDataService.GetHighestPriceBySymbolAndPeriod(r.Context(), symbol, period)
		if err != nil {
			slog.Error("Failed to get highest price by symbol and period", "symbol", symbol, "period", period, "err", err)
			if err == domain.ErrNotFound {
				ResponseError(w, http.StatusNotFound, "Data not found")
			} else {
				ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
			}
			return
		}
		slog.Info("Successfully got highest price by symbol and period", "symbol", symbol, "period", period)
		ResponseJSON(w, http.StatusOK, price)
		return
	}

	price, err := h.marketDataService.GetHighestPriceBySymbol(r.Context(), symbol)
	if err != nil {
		slog.Error("Failed to get highest price by symbol", "symbol", symbol, "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("Successfully got highest price by symbol", "symbol", symbol)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getHighestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.PathValue("exchange")

	period, ok, err := parsePeriod(r)
	if err != nil {
		slog.Error("Invalid period format (highest/exchange)", "symbol", symbol, "exchange", exchange, "err", err)
		ResponseError(w, http.StatusBadRequest, "Invalid period format")
		return
	}

	if ok {
		price, err := h.marketDataService.GetHighestPriceBySymbolAndPeriodAndExchange(r.Context(), symbol, exchange, period)
		if err != nil {
			slog.Error("Failed highest by symbol/period/exchange", "symbol", symbol, "exchange", exchange, "period", period, "err", err)
			if err == domain.ErrNotFound {
				ResponseError(w, http.StatusNotFound, "Data not found")
			} else {
				ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
			}
			return
		}
		slog.Info("OK highest by symbol/period/exchange", "symbol", symbol, "exchange", exchange, "period", period)
		ResponseJSON(w, http.StatusOK, price)
		return
	}

	price, err := h.marketDataService.GetHighestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	if err != nil {
		slog.Error("Failed highest by symbol/exchange", "symbol", symbol, "exchange", exchange, "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("OK highest by symbol/exchange", "symbol", symbol, "exchange", exchange)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getLowestPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	period, ok, err := parsePeriod(r)
	if err != nil {
		slog.Error("Invalid period format (lowest)", "symbol", symbol, "err", err)
		ResponseError(w, http.StatusBadRequest, "Invalid period format")
		return
	}

	if ok {
		price, err := h.marketDataService.GetLowestPriceBySymbolAndPeriod(r.Context(), symbol, period)
		if err != nil {
			slog.Error("Failed lowest by symbol/period", "symbol", symbol, "period", period, "err", err)
			if err == domain.ErrNotFound {
				ResponseError(w, http.StatusNotFound, "Data not found")
			} else {
				ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
			}
			return
		}
		slog.Info("OK lowest by symbol/period", "symbol", symbol, "period", period)
		ResponseJSON(w, http.StatusOK, price)
		return
	}

	price, err := h.marketDataService.GetLowestPriceBySymbol(r.Context(), symbol)
	if err != nil {
		slog.Error("Failed lowest by symbol", "symbol", symbol, "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("OK lowest by symbol", "symbol", symbol)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getLowestPriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.PathValue("exchange")

	period, ok, err := parsePeriod(r)
	if err != nil {
		slog.Error("Invalid period format (lowest/exchange)", "symbol", symbol, "exchange", exchange, "err", err)
		ResponseError(w, http.StatusBadRequest, "Invalid period format")
		return
	}

	if ok {
		price, err := h.marketDataService.GetLowestPriceBySymbolAndPeriodAndExchange(r.Context(), symbol, exchange, period)
		if err != nil {
			slog.Error("Failed lowest by symbol/period/exchange", "symbol", symbol, "exchange", exchange, "period", period, "err", err)
			if err == domain.ErrNotFound {
				ResponseError(w, http.StatusNotFound, "Data not found")
			} else {
				ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
			}
			return
		}
		slog.Info("OK lowest by symbol/period/exchange", "symbol", symbol, "exchange", exchange, "period", period)
		ResponseJSON(w, http.StatusOK, price)
		return
	}

	price, err := h.marketDataService.GetLowestPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	if err != nil {
		slog.Error("Failed lowest by symbol/exchange", "symbol", symbol, "exchange", exchange, "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("OK lowest by symbol/exchange", "symbol", symbol, "exchange", exchange)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getAveragePrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	period, ok, err := parsePeriod(r)
	if err != nil {
		slog.Error("Invalid period format (average)", "symbol", symbol, "err", err)
		ResponseError(w, http.StatusBadRequest, "Invalid period format")
		return
	}

	if ok {
		price, err := h.marketDataService.GetAvgPriceBySymbolAndPeriod(r.Context(), symbol, period)
		if err != nil {
			slog.Error("Failed average by symbol/period", "symbol", symbol, "period", period, "err", err)
			if err == domain.ErrNotFound {
				ResponseError(w, http.StatusNotFound, "Data not found")
			} else {
				ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
			}
			return
		}
		slog.Info("OK average by symbol/period", "symbol", symbol, "period", period)
		ResponseJSON(w, http.StatusOK, price)
		return
	}

	price, err := h.marketDataService.GetAvgPriceBySymbol(r.Context(), symbol)
	if err != nil {
		slog.Error("Failed average by symbol", "symbol", symbol, "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("OK average by symbol", "symbol", symbol)
	ResponseJSON(w, http.StatusOK, price)
}

func (h *Handler) getAveragePriceFromExchange(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	exchange := r.PathValue("exchange")

	period, ok, err := parsePeriod(r)
	if err != nil {
		slog.Error("Invalid period format (average/exchange)", "symbol", symbol, "exchange", exchange, "err", err)
		ResponseError(w, http.StatusBadRequest, "Invalid period format")
		return
	}

	if ok {
		price, err := h.marketDataService.GetAvgPriceBySymbolAndPeriodAndExchange(r.Context(), symbol, exchange, period)
		if err != nil {
			slog.Error("Failed average by symbol/period/exchange", "symbol", symbol, "exchange", exchange, "period", period, "err", err)
			if err == domain.ErrNotFound {
				ResponseError(w, http.StatusNotFound, "Data not found")
			} else {
				ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
			}
			return
		}
		slog.Info("OK average by symbol/period/exchange", "symbol", symbol, "exchange", exchange, "period", period)
		ResponseJSON(w, http.StatusOK, price)
		return
	}

	price, err := h.marketDataService.GetAvgPriceBySymbolAndExchange(r.Context(), symbol, exchange)
	if err != nil {
		slog.Error("Failed average by symbol/exchange", "symbol", symbol, "exchange", exchange, "err", err)
		if err == domain.ErrNotFound {
			ResponseError(w, http.StatusNotFound, "Data not found")
		} else {
			ResponseError(w, http.StatusInternalServerError, "Oops... Something went wrong!")
		}
		return
	}
	slog.Info("OK average by symbol/exchange", "symbol", symbol, "exchange", exchange)
	ResponseJSON(w, http.StatusOK, price)
}
