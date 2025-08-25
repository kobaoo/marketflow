package web

import (
	"marketflow/internal/domain"
	"net/http"
)

type Handler struct {
	marketDataService domain.MarketDataService
}

func NewHandler(marketDataService domain.MarketDataService) *Handler {
	return &Handler{
		marketDataService: marketDataService,
	}
}

func (h *Handler)RegisterRouter(mux *http.ServeMux) {
	// Market Data API
	mux.HandleFunc("GET /prices/latest/{symbol}", h.getLatestPrice)
	mux.HandleFunc("GET /prices/latest/{exchange}/{symbol}", getLatestPriceFromExchange)

	mux.HandleFunc("GET /prices/highest/{symbol}", getHighestPrice)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}", getHighestPriceFromExchange)

	mux.HandleFunc("GET /prices/lowest/{symbol}", getLowestPrice)
	mux.HandleFunc("GET /prices/lowest/{exchange}/{symbol}", getLowestPriceFromExchange)

	mux.HandleFunc("GET /prices/average/{symbol}", getAveragePrice)
	mux.HandleFunc("GET /prices/average/{exchange}/{symbol}", getAveragePriceFromExchange)

	// Data Mode API
	mux.HandleFunc("POST /mode/test", switchToTestMode)
	mux.HandleFunc("POST /mode/live", switchToLiveMode)

	// System Health
	mux.HandleFunc("GET /health", getSystemHealth)
}
