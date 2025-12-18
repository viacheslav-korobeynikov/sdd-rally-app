package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/internal/config"
)

func main() {
	config.Init()
	// Создание инстанса приложения Fiber
	app := fiber.New()

	// Простейший хэндлер
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	// Настраиваем порт, который будет слушать приложение
	app.Listen(":3000")
}
