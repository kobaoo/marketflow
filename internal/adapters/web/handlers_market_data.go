package web

import "net/http"

func getLatestPrice(w http.ResponseWriter, r *http.Request)             {
	
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
