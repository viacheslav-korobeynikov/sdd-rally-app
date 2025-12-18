package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Инициализация загрузки файла окружения
func Init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file")
		return
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

// Структура конфига логов
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

// Структура конфига подключения к БД
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Извлечение конфига БД из переменной окружения
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getString("DB_HOST", "localhost"),
		Port:     getInt("DB_PORT", 5432),
		User:     getString("DB_USER", "postgres"),
		Password: getString("DB_PASSWORD", "postgres"),
		Name:     getString("DB_DATABASE", "rally"),
		SSLMode:  getString("DB_SSLMODE", "disable"),
	}
}
