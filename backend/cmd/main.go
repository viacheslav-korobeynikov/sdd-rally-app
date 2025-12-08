package main

import "github.com/gofiber/fiber/v2"

func main() {

	// Создаем инстанс приложения
	app := fiber.New()
	// Настраиваем порт, который будет слушать приложение
	app.Listen(":3000")

}
