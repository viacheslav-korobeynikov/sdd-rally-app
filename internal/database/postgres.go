package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/config"
)

func CreateDbPool(config *config.DatabaseConfig, logger *zerolog.Logger) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
	)
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error().Msg("Не удалось подключиться к БД")
		panic(err)
	}
	logger.Info().Msg("Подключились к БД")
	return dbpool
}
