package postgres

import (
	"database/sql"
	"log/slog"
	"marketflow/internal/config"
)

func ConnectDB(config *config.Config) (*sql.DB) {
	// Define connection string
	connStr := config.Postgres.Dsn

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
	}
	defer db.Close()

	// Check if the connection works
	err = db.Ping()
	if err != nil {
		slog.Error("Error pinging database", "error", err)
	}

	slog.Info("Connected to Postgres ✅")
	return db
}
