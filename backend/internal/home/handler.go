package home

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/pkg/templadapter"
	"github.com/viacheslav-korobeynikov/sdd-rally-app/views"
)

type HomeHandler struct {
	router       fiber.Router
	customLogger *zerolog.Logger
}

// Функция конструктор
func NewHandler(router fiber.Router, customLogger *zerolog.Logger) {
	h := &HomeHandler{
		router:       router,
		customLogger: customLogger,
	}
	h.router.Get("/", h.home)
}

// Хэндлер для главной страницы
func (h *HomeHandler) home(c *fiber.Ctx) error {
	component := views.Main()
	return templadapter.Render(c, component)
}
