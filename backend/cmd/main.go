package main

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/config"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/home"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/pkg/database"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/pkg/logger"
)

func main() {
	// Получаем переменные окружения
	config.Init()
	// Получаем переменные окружения для БД
	dbConfig := config.NewDatabaseConfig()

	// Получаем переменные окружения для логов
	logConfig := config.NewLogConfig()
	customLogger := logger.NewLogger(logConfig)

	// Создаем инстанс приложения
	app := fiber.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))
	// Middleware, который перезапускает приложение в случае, если произошел вызов panic
	app.Use(recover.New())

	app.Static("/public", "./public") //Обработчик статики (публичных файлов)
	// Создаем подключение к БД
	dbpool := database.CreateDbPool(dbConfig, customLogger)
	defer dbpool.Close()
	home.NewHandler(app, customLogger)
	// Настраиваем порт, который будет слушать приложение
	app.Listen(":3000")

}
