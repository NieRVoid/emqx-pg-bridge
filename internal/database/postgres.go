package database

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/NieRVoid/emqx-pg-bridge/internal/config"
	"github.com/NieRVoid/emqx-pg-bridge/pkg/logger"
)

// Postgres represents a PostgreSQL database connection
type Postgres struct {
	Pool *pgxpool.Pool
	log  *logger.Logger
}

// NewPostgres creates a new PostgreSQL connection pool
func NewPostgres(ctx context.Context, cfg *config.Config, log *logger.Logger) (*Postgres, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	
	// Set pool configuration
	dbConfig.MaxConns = int32(cfg.Database.MaxConnections)
	dbConfig.MinConns = int32(cfg.Database.MinConnections)
	dbConfig.MaxConnLifetime = cfg.GetMaxConnectionLifetime()
	dbConfig.MaxConnIdleTime = cfg.GetMaxConnectionIdleTime()
	
	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	
	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	
	log.Info("Connected to PostgreSQL database")
	
	return &Postgres{
		Pool: pool,
		log:  log,
	}, nil
}

// Close closes the database connection pool
func (p *Postgres) Close() {
	p.Pool.Close()
	p.log.Info("Closed PostgreSQL connection pool")
}