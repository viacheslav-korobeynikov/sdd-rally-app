package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/config"
	"log"
)

func main() {
	// Получаем переменные окружения
	config.Init()
	// Получаем переменные окружения для БД
	dbConfig := config.NewDatabaseConfig()
	log.Println(dbConfig)

	// Создаем инстанс приложения
	app := fiber.New()

	app.Use(recover.New()) // Middleware, который перезапускает приложение в случае, если произошел вызов panic
	// Настраиваем порт, который будет слушать приложение
	app.Listen(":3000")

}
