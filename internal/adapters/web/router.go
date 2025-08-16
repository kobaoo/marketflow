package web

import (
	"net/http"
)

func RegisterRouter() {
	mux := http.NewServeMux()

	// Market Data API
	mux.HandleFunc("/prices/latest/{symbol}", getLatestPrice)
	mux.HandleFunc("/prices/latest/{exchange}/{symbol}", getLatestPriceFromExchange)
	mux.HandleFunc("/prices/highest/{symbol}", getHighestPrice)
	mux.HandleFunc("/prices/highest/{exchange}/{symbol}", getHighestPriceFromExchange)
	mux.HandleFunc("/prices/highest/{symbol}", getHighestPriceWithPeriod)
	mux.HandleFunc("/prices/highest/{exchange}/{symbol}", getHighestPriceFromExchangeWithPeriod)
	mux.HandleFunc("/prices/lowest/{symbol}", getLowestPrice)
	mux.HandleFunc("/prices/lowest/{exchange}/{symbol}", getLowestPriceFromExchange)
	mux.HandleFunc("/prices/lowest/{symbol}", getLowestPriceWithPeriod)
	mux.HandleFunc("/prices/lowest/{exchange}/{symbol}", getLowestPriceFromExchangeWithPeriod)
	mux.HandleFunc("/prices/average/{symbol}", getAveragePrice)
	mux.HandleFunc("/prices/average/{exchange}/{symbol}", getAveragePriceFromExchange)
	mux.HandleFunc("/prices/average/{exchange}/{symbol}", getAveragePriceFromExchangeWithPeriod)

	// Data Mode API
	mux.HandleFunc("/mode/test", switchToTestMode)
	mux.HandleFunc("/mode/live", switchToLiveMode)

	// System Health
	mux.HandleFunc("/health", getSystemHealth)
}
