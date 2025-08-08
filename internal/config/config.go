package config

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServerConfig *ServerConfig
}

type ServerConfig struct {
	Port int
}

func NewConfig() *Config {
	var serverConfig ServerConfig
	flag.IntVar(&serverConfig.Port, "port", 0, "Port to serve on")
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
