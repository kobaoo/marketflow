package web

import (
	"log/slog"
	"marketflow/internal/config"
	"net/http"
)

func StartServer(config config.Config) error {
	mux := http.NewServeMux()
	RegisterRouter(mux)
	
	server := &http.Server{Addr: ":" + config.Port, Handler: mux}

	slog.Info("Server is starting", "port", config.Port)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
