package cmd

import (
	// "log/slog"
	// "marketflow/internal/config"

	"log/slog"
	"marketflow/internal/adapters/exchange"
	"marketflow/internal/infra"
)

func RunApp() {
	infra.SetUpLogger()
	allMessages := make(chan []byte, 30)
	go exchange.RunTCPClients(allMessages, false)
	//just for check
	for {
		select{ 
		case msg, ok := <-allMessages: 
			if !ok {
				break
			}
			slog.Debug("GOT message", "msg", msg)
		}
	}

	// cfg, err := config.Load("./configs/config.json")
	// if err != nil {
	// 	slog.Error("Config Error", "error", err)
	// }
	
}
