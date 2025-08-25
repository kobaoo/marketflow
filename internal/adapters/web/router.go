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
	mux.HandleFunc("GET /prices/latest/{exchange}/{symbol}", h.getLatestPriceFromExchange)

	mux.HandleFunc("GET /prices/highest/{symbol}", h.getHighestPrice)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}", h.getHighestPriceFromExchange)

	mux.HandleFunc("GET /prices/lowest/{symbol}", h.getLowestPrice)
	mux.HandleFunc("GET /prices/lowest/{exchange}/{symbol}", h.getLowestPriceFromExchange)

	mux.HandleFunc("GET /prices/average/{symbol}", h.getAveragePrice)
	mux.HandleFunc("GET /prices/average/{exchange}/{symbol}", h.getAveragePriceFromExchange)

	// Data Mode API
	mux.HandleFunc("POST /mode/test", switchToTestMode)
	mux.HandleFunc("POST /mode/live", switchToLiveMode)

	// System Health
	mux.HandleFunc("GET /health", getSystemHealth)
}
