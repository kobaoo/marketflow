package web

import "net/http"

func getLatestPrice(w http.ResponseWriter, r *http.Request)             {}
func getLatestPriceFromExchange(w http.ResponseWriter, r *http.Request) {}

func getHighestPrice(w http.ResponseWriter, r *http.Request)                       {}
func getHighestPriceFromExchange(w http.ResponseWriter, r *http.Request)           {}

func getLowestPrice(w http.ResponseWriter, r *http.Request)                       {}
func getLowestPriceFromExchange(w http.ResponseWriter, r *http.Request)           {}

func getAveragePrice(w http.ResponseWriter, r *http.Request)                       {}
func getAveragePriceFromExchange(w http.ResponseWriter, r *http.Request)           {}