package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	HTTP HTTP `json:"http"`
	Mode string `json:"mode"`
	Pairs []string `json:"pairs"`
	Exchanges []Exchange `json:"exchanges"`
	AggWindowS int `agg_window_seconds"`
}

type Exchange struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type HTTP struct {
	Port    int `json:"port"`
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	if c.HTTP.Port < 1024 || c.HTTP.Port > 65565 {
		slog.Warn("Invalid or empty port number. Loading default :8080")
		c.HTTP.Port = 8080
	}
	if c.AggWindowS <= 0 {
		slog.Warn("Invalid or empty window seconds for aggregation. Loading default 60s")
		c.AggWindowS = 60
	}
	if (c.Mode == "" || c.Mode != "live") && c.Mode != "test"{
		slog.Warn("Invalid or empty mode definition. Loading default <test>")
		c.Mode = "test"
	}
	return c, nil
}

// func NewConfig() (*Config, error) {
// 	var serverConfig ServerConfig
// 	flag.IntVar(&serverConfig.Port, "port", 0, "Port to serve on")
// 	flag.Usage = printHelp
// 	flag.Parse()

// 	if serverConfig.Port == 0 {
// 		strPort := getEnv("SERVER_PORT", "8080")
// 		port, err := strconv.Atoi(strPort)
// 		if err == nil {
// 			serverConfig.Port = port
// 		} else {
// 			slog.Error("Invalid port number in environment, setting default :8080")
// 			serverConfig.Port = 8080
// 		}
// 	}

// 	serverConfig.CfgPath = getEnv("SERVER_CONFIG_PATH", "configs/config.json")
// 	return &Config{
// 		ServerConfig: &serverConfig,
// 	}
// }

// func getEnv(name, defaultValue string) string {
// 	value := os.Getenv(name)
// 	if value == "" {
// 		slog.Warn("Missing environment variables, setting default values!", "env", name, "default value", defaultValue)
// 		value = defaultValue
// 	}
// 	return value
// }

func printHelp() {
	fmt.Println(`
Usage:
  marketflow [--port <N>]
  marketflow --help

Options:
  --port N     Port number`)
}
