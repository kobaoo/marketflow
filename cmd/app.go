package cmd

import (
	// "log/slog"
	// "marketflow/internal/config"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/infra"
)

func RunApp() {
	infra.SetUpLogger()

	exchange.RunTCPClients()

	// cfg, err := config.Load("./configs/config.json")
	// if err != nil {
	// 	slog.Error("Config Error", "error", err)
	// }
	
}
