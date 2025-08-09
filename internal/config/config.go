package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServerConfig *ServerConfig
}

type ServerConfig struct {
	Port    int `json:"port"`
	CfgPath string
}

func NewConfig() *Config {
	var serverConfig ServerConfig
	flag.IntVar(&serverConfig.Port, "port", 0, "Port to serve on")
	flag.Usage = printHelp
	flag.Parse()

	if serverConfig.Port == 0 {
		strPort := getEnv("SERVER_PORT", "8080")
		port, err := strconv.Atoi(strPort)
		if err == nil {
			serverConfig.Port = port
		} else {
			slog.Error("Invalid port number in environment, setting default :8080")
			serverConfig.Port = 8080
		}
	}

	serverConfig.CfgPath = getEnv("SERVER_CONFIG_PATH", "configs/config.json")
	return &Config{
		ServerConfig: &serverConfig,
	}
}

func getEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		slog.Warn("Missing environment variables, setting default values!", "env", name, "default value", defaultValue)
		value = defaultValue
	}
	return value
}

func printHelp() {
	fmt.Println(`
Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number`)
}
