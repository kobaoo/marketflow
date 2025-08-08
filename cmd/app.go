package cmd

import (
	"marketflow/internal/config"
	"marketflow/internal/logger"
)

func RunApp() {
	logger.SetLogger()

	config := config.NewConfig()
}
