package cmd

import (
	"marketflow/internal/config"
	"marketflow/internal/infra"
)

func RunApp() {
	infra.SetLogger()

	config := config.NewConfig()
}
