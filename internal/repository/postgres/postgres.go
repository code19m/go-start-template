package postgres

import (
	"context"
	"fmt"
	"go-start-template/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func NewConnPool(cfg *config.Postgres) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pool, nil
}

// default configuration of connection pool
// var defaultMaxConns = int32(4)
// var defaultMinConns = int32(0)
// var defaultMaxConnLifetime = time.Hour
// var defaultMaxConnIdleTime = time.Minute * 30
// var defaultHealthCheckPeriod = time.Minute
