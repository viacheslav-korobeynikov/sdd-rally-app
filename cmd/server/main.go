package main

import "github.com/gofiber/fiber/v2"

func main() {
	// Создание инстанса приложения Fiber
	app := fiber.New()

	// Простейший хэндлер
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	// Настраиваем порт, который будет слушать приложение
	app.Listen(":3000")
}
