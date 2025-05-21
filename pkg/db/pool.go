package db

import (
	"GinBase/pkg/config"
	"GinBase/pkg/logger"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

var (
	poolOnce sync.Once
)

type Pool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Close()
}

type PostgresPool struct {
	pool *pgxpool.Pool
	log  logger.Logger
}

func NewPostgresPool(config *config.DatabaseConfig, log logger.Logger) (*PostgresPool, error) {
	if config == nil {
		return nil, config.ErrConfigNotFound
	}
}
