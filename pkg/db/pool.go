package db

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"sync"
	"terraqt.io/colas/bedrock-go/pkg/config"
	"terraqt.io/colas/bedrock-go/pkg/errs"
	"terraqt.io/colas/bedrock-go/pkg/logger"
	"time"
)

var (
	poolOnce sync.Once
	pool     PGPool
	poolErr  error
)

type PGPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	driver.Pinger
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

type adaptPool interface {
	PGPool
	getStdPool() *pgxpool.Pool
}

type contextKey string

const queryStartTimeKey contextKey = "queryStartTime"

type sqlTracer struct {
	logger.Logger
}

func (s *sqlTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	startTime := time.Now()
	s.Debug(ctx, "SQL execution started", zap.String("sql", data.SQL))
	return context.WithValue(ctx, queryStartTimeKey, startTime)
}

func (s *sqlTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	startTime, ok := ctx.Value(queryStartTimeKey).(time.Time)
	if !ok {
		s.Warn(ctx, "Failed to get query start time from context")
		return
	}

	if data.Err != nil {
		s.Error(ctx, "SQL execution failed", zap.Error(data.Err))
		return
	}

	duration := time.Since(startTime)
	s.Debug(ctx, "SQL execution completed", zap.Duration("duration", duration))
}

type postgresPool struct {
	*pgxpool.Pool
	log logger.Logger
}

func (p *postgresPool) getStdPool() *pgxpool.Pool {
	return p.Pool
}

func newPostgresPool(config config.DatabaseConfig, log logger.Logger) (PGPool, error) {

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
	)

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error(nil, "failed to parse postgres connection string", zap.Error(err))
		return nil, errs.WrapCodeError(
			errs.ErrInvalidParam,
			fmt.Errorf("failed to parse postgres connection string: %w", err),
		)
	}

	// 设置连接池参数
	pgxConfig.MaxConns = config.MaxOpenConns
	pgxConfig.MinConns = config.MaxIdleConns
	pgxConfig.MaxConnLifetime = config.ConnMaxLifetime
	pgxConfig.MaxConnIdleTime = config.ConnMaxIdleTime

	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	if config.DebugSQL {
		pgxConfig.ConnConfig.Tracer = &sqlTracer{log}
		log.Info(nil, "SQL debug mode is enabled")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Error(nil, "failed to create postgres connection pool", zap.Error(err))
		return nil, errs.WrapCodeError(
			errs.ErrResourceInitFailed,
			fmt.Errorf("failed to create postgres connection pool: %w", err),
		)
	}

	// 验证连接池是否可以正常连接
	if err := dbPool.Ping(ctx); err != nil {
		log.Error(nil, "failed to ping postgres database", zap.Error(err))
		dbPool.Close()
		return nil, errs.WrapCodeError(
			errs.ErrDBConnection,
			fmt.Errorf("failed to ping postgres database: %w", err),
		)
	}

	log.Info(nil, "Successfully connected to database")

	return &postgresPool{
		Pool: dbPool,
		log:  log,
	}, nil
}

func providePostgresPool(config config.DatabaseConfig, log logger.Logger) (PGPool, error) {
	poolOnce.Do(
		func() {
			pool, poolErr = newPostgresPool(config, log)
		},
	)

	if poolErr != nil {
		log.Error(nil, "failed to create postgres connection pool", zap.Error(poolErr))
		return nil, errs.WrapCodeError(
			errs.ErrResourceInitFailed,
			fmt.Errorf("failed to create postgres connection pool: %w", poolErr),
		)
	}
	return pool, poolErr
}

// Acquire 获取数据库连接
func (p *postgresPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := p.Pool.Acquire(ctx)
	if err != nil {
		p.log.Error(nil, "failed to acquire database connection", zap.Error(err))
		return nil, errs.WrapCodeError(
			errs.ErrDBConnection,
			fmt.Errorf("failed to acquire database connection: %w", err),
		)
	}
	return conn, nil
}

func (p *postgresPool) Ping(ctx context.Context) error {
	if err := p.Pool.Ping(ctx); err != nil {
		p.log.Error(nil, "failed to ping database", zap.Error(err))
		return errs.WrapCodeError(
			errs.ErrDBConnection,
			fmt.Errorf("failed to ping database: %w", err),
		)
	}
	return nil
}

// BeginTx starts a new database transaction using the provided context and transaction options.
// If options are nil, the default transaction ReadCommitted is applied.
// Returns the transaction object or an error if the transaction initialization fails.
func (p *postgresPool) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {

	tx, err := p.Pool.BeginTx(ctx, opts)

	if err != nil {
		p.log.Error(nil, "failed to begin transaction", zap.Error(err))
		return nil, errs.WrapCodeError(
			errs.ErrDBTransaction,
			fmt.Errorf("failed to begin transaction: %w", err),
		)
	}

	return tx, nil
}
