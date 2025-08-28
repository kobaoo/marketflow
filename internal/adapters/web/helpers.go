package web

import (
	"net/http"
	"time"
)

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