package web

import (
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) switchToTestMode(w http.ResponseWriter, r *http.Request) {
	err := h.systemService.SwitchToTestMode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	currentMode, _ := h.systemService.GetCurrentMode(r.Context())
	response := map[string]interface{}{
		"status":  "success",
		"message": "Switched to test mode",
		"mode":    currentMode,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) switchToLiveMode(w http.ResponseWriter, r *http.Request) {
	err := h.systemService.SwitchToLiveMode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	currentMode, _ := h.systemService.GetCurrentMode(r.Context())
	response := map[string]interface{}{
		"status":  "success",
		"message": "Switched to live mode",
		"mode":    currentMode,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) getSystemHealth(w http.ResponseWriter, r *http.Request) {
	health, err := h.systemService.GetSystemHealth(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := map[string]interface{}{
		"status":      health.Status,
		"service":     health.Service,
		"message":     health.Message,
		"redis_ok":    health.Redis,
		"postgres_ok": health.PostgreSQL,
		"sources":     health.Exchanges, 
		"mode":        health.Mode,
		"timestamp":   health.Timestamp.Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}