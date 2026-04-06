package main

import (
	"fmt"
	"log"

	"github.com/financas/backend/internal/config"
	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/handlers"
	"github.com/financas/backend/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Carregar configurações
	log.Println("Carregando configurações...")
	config.LoadConfig()

	// Conectar ao banco de dados
	log.Println("Conectando ao banco de dados...")
	if err := database.Connect(); err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Inicializar serviços
	log.Println("Inicializando serviços de IA...")
	handlers.InitAIService(database.DB)

	// Criar aplicação Fiber
	app := fiber.New(fiber.Config{
		AppName:      "Financas API v1.0",
		ErrorHandler: customErrorHandler,
	})

	// Middlewares globais
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "API Financas rodando",
		})
	})

	// Configurar rotas
	routes.SetupRoutes(app)

	// Iniciar servidor
	port := config.AppConfig.Port
	log.Printf("🚀 Servidor rodando na porta %s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}

// customErrorHandler trata erros globalmente
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
