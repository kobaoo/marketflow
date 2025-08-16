package web

import "net/http"

func switchToTestMode(w http.ResponseWriter, r *http.Request) {}
func switchToLiveMode(w http.ResponseWriter, r *http.Request) {}
func getSystemHealth(w http.ResponseWriter, r *http.Request)  {}
