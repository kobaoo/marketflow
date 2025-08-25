package web

import (
	"context"
	"log/slog"
	"marketflow/internal/config"
	"net/http"
)

func (h *Handler) StartServer(ctx context.Context, config *config.Config) error {
	mux := http.NewServeMux()
	h.RegisterRouter(mux)

	server := &http.Server{Addr: ":" + config.Port, Handler: mux}

	slog.Info("Server is starting", "port", config.Port)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server startup error", "error", err)
		}
	}()

	slog.Info("Server is listening", "port", config.Port, "address", "http://localhost:"+config.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	slog.Info("Server is shutting down...")

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown error", "error", err)
		return err
	}

	slog.Info("Server stopped")
	return nil
}
