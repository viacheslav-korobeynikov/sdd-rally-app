package middleware

import "github.com/gofiber/fiber/v2"

// SecurityHeadersConfig конфигурация security headers
type SecurityHeadersConfig struct {
	CSPDefaultSrc string
}

// ConfigDefault настройки по умолчанию
var ConfigDefault = SecurityHeadersConfig{
	CSPDefaultSrc: "default-src 'self'; frame-ancestors 'self'; base-uri 'self'; form-action 'self'",
}

// Создаем middleware для security headers
func NewSecurityHeaders(config SecurityHeadersConfig) fiber.Handler {
	cfg := ConfigDefault
	if config.CSPDefaultSrc != "" {
		cfg.CSPDefaultSrc = config.CSPDefaultSrc
	}
	return func(c *fiber.Ctx) error {
		// Устанавливаем security headers
		c.Set("Content-Security-Policy", cfg.CSPDefaultSrc)
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-Content-Type-Options", "nosniff")
		return c.Next()
	}
}
