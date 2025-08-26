package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"marketflow/internal/config"
	"time"
)


func ConnectDB(cfg *config.Config) *sql.DB {
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.Postgres.User,
        cfg.Postgres.Password,
        cfg.Postgres.Host,
        cfg.Postgres.Port,
        cfg.Postgres.DBName,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        slog.Error("Error connecting to database (open)", "error", err)
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        slog.Error("Error pinging database", "error", err)
        return nil
    }

    slog.Info("Connected to Postgres ✅")
    return db
}