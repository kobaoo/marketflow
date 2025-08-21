package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func ConnectDB() {
	// Define connection string
	connStr := "postgres://user:password@localhost:5432/mydb?sslmode=disable"

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

	fmt.Println("Connected to Postgres ✅")
}
