package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/config"
)

// Кастомный логгер
func NewLogger(config *config.LogConfig) *zerolog.Logger {
	//Устанавливаем уровень логов в заисимости от окружения
	zerolog.SetGlobalLevel(zerolog.Level(config.Level))
	var logger zerolog.Logger
	if config.Format == "json" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		consoleWritter := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(consoleWritter).With().Timestamp().Logger()
	}
	return &logger
}
