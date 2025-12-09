package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/config"
)

// Функция для создания подключения к БД
func CreateDbPool(config *config.DatabaseConfig, logger *zerolog.Logger) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), config.Url)
	if err != nil {
		logger.Error().Msg("Не удалось подключиться к БД")
		panic(err)
	}
	logger.Info().Msg("Подключились к БД")
	return dbpool
}
