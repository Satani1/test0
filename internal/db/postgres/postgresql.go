package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type DB struct {
}

func NewClient(ctx context.Context, maxAttempts int, username, password, host, port, database string) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:/%s/%s", username, password, host, port, database)

	for maxAttempts > 0 {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err := pgxpool.Connect(ctx, dsn)
		if err != nil {
			fmt.Print("failed to connect to postgresql")
			return
		}
	}

	return nil

}
