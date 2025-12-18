package main

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/config"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/database"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/middleware"
)

func main() {
	// Получение данных из файла конфигурации
	config.Init()
	// Вызов конфигурации логов
	logConfig := config.NewLogConfig()
	// Вызов конфигурации БД
	dbConfig := config.NewDatabaseConfig()
	customLogger := middleware.NewLogger(logConfig)
	// Создание инстанса приложения Fiber
	app := fiber.New()

	app.Use(recover.New())

	dbpool := database.CreateDbPool(dbConfig, customLogger)
	defer dbpool.Close()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: customLogger,
	}))

	// Простейший хэндлер
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	// Настраиваем порт, который будет слушать приложение
	app.Listen(":3000")
}
