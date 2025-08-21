package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port      int            `json:"port"`
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
	Dsn string `json:"dsn"`
}

func loadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
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
