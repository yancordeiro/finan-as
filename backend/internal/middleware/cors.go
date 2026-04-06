package middleware

import (
	"github.com/financas/backend/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS configura o middleware de CORS
func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.FrontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           3600,
	})
}
