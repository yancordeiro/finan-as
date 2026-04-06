package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/financas/backend/internal/middleware"
	"github.com/financas/backend/internal/services"
	"github.com/financas/backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// GetDashboardSummary retorna o resumo financeiro do mês
func GetDashboardSummary(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse dos parâmetros de mês/ano (padrão: mês atual)
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil && y >= 2000 && y <= 2100 {
			year = y
		}
	}

	summary, err := services.GetDashboardSummary(userID, month, year)
	if err != nil {
		log.Printf("Erro ao buscar resumo do dashboard: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar resumo")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Resumo encontrado", summary)
}

// GetExpensesByCategory retorna gastos agrupados por categoria
func GetExpensesByCategory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil && y >= 2000 && y <= 2100 {
			year = y
		}
	}

	expenses, err := services.GetExpensesByCategory(userID, month, year)
	if err != nil {
		log.Printf("Erro ao buscar gastos por categoria: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar gastos por categoria")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Gastos por categoria encontrados", expenses)
}

// GetIncomesByCategory retorna entradas agrupadas por categoria
func GetIncomesByCategory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil && y >= 2000 && y <= 2100 {
			year = y
		}
	}

	incomes, err := services.GetIncomesByCategory(userID, month, year)
	if err != nil {
		log.Printf("Erro ao buscar entradas por categoria: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar entradas por categoria")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Entradas por categoria encontradas", incomes)
}

// GetMonthlyFlow retorna o fluxo financeiro dos últimos meses
func GetMonthlyFlow(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Padrão: últimos 6 meses
	months := 6
	if monthsStr := c.Query("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m >= 1 && m <= 12 {
			months = m
		}
	}

	flow, err := services.GetMonthlyFlow(userID, months)
	if err != nil {
		log.Printf("Erro ao buscar fluxo mensal: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar fluxo mensal")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Fluxo mensal encontrado", flow)
}

// GetRecentTransactions retorna as últimas transações
func GetRecentTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Padrão: 10 transações
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l >= 1 && l <= 50 {
			limit = l
		}
	}

	transactions, err := services.GetRecentTransactions(userID, limit)
	if err != nil {
		log.Printf("Erro ao buscar transações recentes: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar transações recentes")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transações recentes encontradas", transactions)
}
