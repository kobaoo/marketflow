package cmd

import (
	"log/slog"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/config"
	"marketflow/internal/infra"
)

func RunApp() {
	infra.SetUpLogger()

	config, err := config.ReadConfig() // Read config
	if err != nil {
		slog.Error("Config Error", "error", err)
	}

	exchange.RunTCPClients(false)
}
