package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port      string         `json:"port"`
	Mode      string         `json:"mode"`
	Exchanges []Exchange     `json:"exchanges"`
	Redis     RedisConfig    `json:"redis"`
	Postgres  PostgresConfig `json:"postgres"`
}

type Exchange struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func loadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return config, err
	}

	if Port != 0 {
		port = Port
	}

	if port <= 1024 || port > 65565 {
		return config, fmt.Errorf("invalid port number, should be between 1024 and 65565")
	}

	if config.Mode != "test" && config.Mode != "live" {
		return config, fmt.Errorf("invalid mode, expected live or test")
	}

	if config.Postgres.Port <= 1024 || config.Postgres.Port > 65565 {
		return config, fmt.Errorf("invalid port number for postgres, should be between 1024 and 65565")
	}

	return config, err
}

func ReadConfig() (Config, error) {
	configPath := "configs/config.json"
	config, err := loadConfig(configPath)
	if err != nil {
		return config, err
	}
	return config, nil
}
