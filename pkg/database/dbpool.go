package database

import (
	"context"
	"fmt"
	"jwt-service/config"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	cfg, err := config.LoadCfg()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Server.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLmode,
	)

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Error parsing database DSN: ", err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
	if err != nil {
		wrappedErr := fmt.Errorf("error connecting to database: %w", err)
		log.WithError(wrappedErr).Error("Init db")
		return nil, wrappedErr
	}

	return pool, nil
}
