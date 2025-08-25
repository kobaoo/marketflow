package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port      string            `json:"port"`
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

	// Override with environment variables for Docker
	applyEnvOverrides(&config)
	return config, nil
}

func applyEnvOverrides(config *Config) {
	// Override Redis configuration
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		redisPort := os.Getenv("REDIS_PORT")
		if redisPort == "" {
			redisPort = "6379"
		}
		config.Redis.Addr = fmt.Sprintf("%s:%s", redisHost, redisPort)
	}

	// Override PostgreSQL configuration
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		if dbUser == "" {
			dbUser = "postgres"
		}
		if dbPass == "" {
			dbPass = "secret"
		}
		config.Postgres.Dsn = fmt.Sprintf("postgres://%s:%s@%s:5432/mydb?sslmode=disable", dbUser, dbPass, dbHost)
	}

	// Override exchange configurations
	if ex1Host := os.Getenv("EXCHANGE1_HOST"); ex1Host != "" {
		for i := range config.Exchanges {
			if config.Exchanges[i].Name == "ex1" {
				config.Exchanges[i].Addr = fmt.Sprintf("%s:40101", ex1Host)
			}
		}
	}
	if ex2Host := os.Getenv("EXCHANGE2_HOST"); ex2Host != "" {
		for i := range config.Exchanges {
			if config.Exchanges[i].Name == "ex2" {
				config.Exchanges[i].Addr = fmt.Sprintf("%s:40102", ex2Host)
			}
		}
	}
	if ex3Host := os.Getenv("EXCHANGE3_HOST"); ex3Host != "" {
		for i := range config.Exchanges {
			if config.Exchanges[i].Name == "ex3" {
				config.Exchanges[i].Addr = fmt.Sprintf("%s:40103", ex3Host)
			}
		}
	}
}
