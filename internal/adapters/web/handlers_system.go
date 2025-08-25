package web

import (
	"encoding/json"
	"net/http"
)

func switchToTestMode(w http.ResponseWriter, r *http.Request) {}
func switchToLiveMode(w http.ResponseWriter, r *http.Request) {}

func getSystemHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"service": "marketflow",
		"message": "Service is running",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
