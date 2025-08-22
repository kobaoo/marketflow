package web

import (
	"net/http"
)

func RegisterRouter(mux *http.ServeMux) {
	// Market Data API
	mux.HandleFunc("GET /prices/latest/{symbol}", getLatestPrice)
	mux.HandleFunc("GET /prices/latest/{exchange}/{symbol}", getLatestPriceFromExchange)

	mux.HandleFunc("GET /prices/highest/{symbol}", getHighestPrice)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}", getHighestPriceFromExchange)
	mux.HandleFunc("GET /prices/highest/{symbol}?period={duration}", getHighestPriceWithPeriod)
	mux.HandleFunc("GET /prices/highest/{exchange}/{symbol}?period={duration}", getHighestPriceFromExchangeWithPeriod)

	mux.HandleFunc("GET /prices/lowest/{symbol}", getLowestPrice)
	mux.HandleFunc("GET /prices/lowest/{exchange}/{symbol}", getLowestPriceFromExchange)
	mux.HandleFunc("GET /prices/lowest/{symbol}?period={duration}", getLowestPriceWithPeriod)
	mux.HandleFunc("GET /prices/lowest/{exchange}/{symbol}?period={duration}", getLowestPriceFromExchangeWithPeriod)

	mux.HandleFunc("GET /prices/average/{symbol}", getAveragePrice)
	mux.HandleFunc("GET /prices/average/{exchange}/{symbol}", getAveragePriceFromExchange)
	mux.HandleFunc("GET /prices/average/{symbol}?period={duration}", getAveragePriceWithPeriod)
	mux.HandleFunc("GET /prices/average/{exchange}/{symbol}?period={duration}", getAveragePriceFromExchangeWithPeriod)

	// Data Mode API
	mux.HandleFunc("POST /mode/test", switchToTestMode)
	mux.HandleFunc("POST /mode/live", switchToLiveMode)

	// System Health
	mux.HandleFunc("GET /health", getSystemHealth)
}
