package middleware

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/config"
)

// Кастомный логгер
func NewLogger(config *config.LogConfig) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.Level(config.Level))
	var logger zerolog.Logger
	if config.Format == "json" {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	}
	return &logger
}
