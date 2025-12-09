package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Инициализация загрузки файла окружения
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: No .env file")
	}
	log.Println(".env file loaded")
}

// Работа со строковыми значениями
func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Работа с числовыми значениями
func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

// Работа с булево значениями
func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

// Структура конфигурации БД
type DatabaseConfig struct {
	Url string
}

// Извление конфига БД из переменной окружения
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Url: getString("POSTGRES_DSN", ""),
	}
}

// Структура конфигурации логов
type LogConfig struct {
	Level  int
	Format string
}

// Извлечение конфига логов из переменной окружения
func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  getInt("LOG_LEVEL", 0),
		Format: getString("LOG_FORMAT", "json"),
	}
}

// Структура SMTP
type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Извленение данных для SMTP из переменных окружения
func NewSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:     getString("SMTP_HOST", "mailhog"),
		Port:     getString("SMTP_PORT", "1025"),
		User:     getString("SMTP_USER", ""),
		Password: getString("SMTP_PASSWORD", ""),
	}
}
