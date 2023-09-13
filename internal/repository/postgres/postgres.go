package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type Config struct {
	Host         string
	Port         uint16
	Username     string
	Name         string
	Password     string
	SSLmode      string
	PoolMaxConns uint16
}

func InitDB(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dbCfg, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s host=%s port=%d"+
		" dbname=%s sslmode=%s pool_max_conns=%d", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLmode, cfg.PoolMaxConns))
	if err != nil {
		return nil, fmt.Errorf("parse config to database: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return pool, nil
}
