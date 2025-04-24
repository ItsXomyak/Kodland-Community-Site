package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, dbURL string) (*pgx.Conn, error) {
	const maxAttempts = 10
	const retryDelay = 2 * time.Second

	var conn *pgx.Conn
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err = pgx.Connect(ctx, dbURL)
		if err == nil {
			if err = conn.Ping(ctx); err == nil {
				return conn, nil
			}
			conn.Close(ctx)
		}
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
}