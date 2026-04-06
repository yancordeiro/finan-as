package routes

import (
	"github.com/financas/backend/internal/handlers"
	"github.com/financas/backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas as rotas da aplicação
func SetupRoutes(app *fiber.App) {
	// Aplicar CORS globalmente
	app.Use(middleware.SetupCORS())

	// Grupo de rotas da API
	api := app.Group("/api")

	// Rotas de autenticação (públicas)
	auth := api.Group("/auth")
	{
		auth.Post("/google", handlers.GoogleAuth)
		auth.Post("/refresh", handlers.RefreshToken)
		auth.Post("/logout", middleware.AuthMiddleware, handlers.Logout)
	}

	// Rotas protegidas (requerem autenticação)
	protected := api.Group("", middleware.AuthMiddleware)

	// User
	protected.Get("/user/me", handlers.GetCurrentUser)

	// Categorias
	categories := protected.Group("/categories")
	{
		categories.Get("/", handlers.GetCategories)
		categories.Get("/:id", handlers.GetCategory)
		categories.Post("/", handlers.CreateCategory)
		categories.Put("/:id", handlers.UpdateCategory)
		categories.Delete("/:id", handlers.DeleteCategory)
	}

	// Transações
	transactions := protected.Group("/transactions")
	{
		transactions.Get("/", handlers.GetTransactions)
		transactions.Get("/export", handlers.ExportTransactions)
		transactions.Get("/:id", handlers.GetTransaction)
		transactions.Post("/", handlers.CreateTransaction)
		transactions.Put("/:id", handlers.UpdateTransaction)
		transactions.Delete("/:id", handlers.DeleteTransaction)
	}

	// Dashboard
	dashboard := protected.Group("/dashboard")
	{
		dashboard.Get("/summary", handlers.GetDashboardSummary)
		dashboard.Get("/expenses-by-category", handlers.GetExpensesByCategory)
		dashboard.Get("/incomes-by-category", handlers.GetIncomesByCategory)
		dashboard.Get("/monthly-flow", handlers.GetMonthlyFlow)
		dashboard.Get("/recent-transactions", handlers.GetRecentTransactions)
	}

	// OFX Import
	ofx := protected.Group("/ofx")
	{
		ofx.Post("/preview", handlers.PreviewOFX)
		ofx.Post("/import", handlers.ImportOFX)
	}

	// AI Assistant
	ai := protected.Group("/ai")
	{
		ai.Post("/chat", handlers.ChatWithAI)
		ai.Get("/summary", handlers.GetFinancialSummary)
		ai.Get("/quick-analysis", handlers.QuickAnalysis)
	}

	// Rota de teste para verificar autenticação
	protected.Get("/test", func(c *fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		return c.JSON(fiber.Map{
			"message": "Protected route works",
			"userId":  userID,
		})
	})
}
