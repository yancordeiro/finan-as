package handlers

import (
	"github.com/financas/backend/internal/middleware"
	"github.com/financas/backend/internal/models"
	"github.com/financas/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var aiService *services.AIService

// InitAIService inicializa o serviço de IA
func InitAIService(db *gorm.DB) {
	aiService = services.NewAIService(db)
}

// ChatWithAI processa mensagem do usuário e retorna resposta da IA
// POST /api/ai/chat
func ChatWithAI(c *fiber.Ctx) error {
	// Pegar userID do contexto (autenticação)
	userID := middleware.GetUserID(c)

	// Parse request
	var req models.AIRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validar mensagem
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message is required",
		})
	}

	if len(req.Message) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message too long (max 1000 characters)",
		})
	}

	// Chamar serviço de IA
	response, err := aiService.Chat(userID, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process AI request",
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}

// GetFinancialSummary retorna resumo financeiro do usuário
// GET /api/ai/summary
func GetFinancialSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	summary, err := aiService.GetFinancialSummary(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get financial summary",
		})
	}

	return c.JSON(summary)
}

// QuickAnalysis retorna análises rápidas pré-definidas
// GET /api/ai/quick-analysis
func QuickAnalysis(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Análises pré-definidas
	analyses := []string{
		"Analise meus gastos do mês e identifique onde posso economizar",
		"Quais são minhas maiores despesas e como reduzi-las?",
		"Estou gastando muito? Sugira um plano de economia",
		"Por que meu saldo está negativo este mês?",
	}

	return c.JSON(fiber.Map{
		"userId":   userID,
		"analyses": analyses,
	})
}
