package main

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/config"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/middleware"
)

func main() {
	config.Init()
	logConfig := config.NewLogConfig()
	customLogger := middleware.NewLogger(logConfig)
	// Создание инстанса приложения Fiber
	app := fiber.New()

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
