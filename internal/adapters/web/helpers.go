package web

import (
	"encoding/json"
	"net/http"
	"time"
)

func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	ResponseJSON(w, status, map[string]string{
		"error": message,
	})
}

func parsePeriod(r *http.Request) (time.Duration, bool, error) {
	q := r.URL.Query().Get("period")
	if q == "" {
		return 0, false, nil
	}

	d, err := time.ParseDuration(q)
	if err != nil {
		return 0, false, err
	}

	return d, true, nil
}
