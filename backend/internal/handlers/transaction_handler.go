package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/financas/backend/internal/middleware"
	"github.com/financas/backend/internal/models"
	"github.com/financas/backend/internal/services"
	"github.com/financas/backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// TransactionRequest representa o corpo da requisição para criar/atualizar transação
type TransactionRequest struct {
	Type        models.TransactionType `json:"type"`
	Description string                 `json:"description"`
	Amount      int64                  `json:"amount"` // Em centavos
	Date        string                 `json:"date"`   // Formato: YYYY-MM-DD
	CategoryID  uint                   `json:"categoryId"`
	Notes       string                 `json:"notes"`
}

// GetTransactions retorna as transações do usuário com filtros e paginação
func GetTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse dos filtros da query string
	filters := services.TransactionFilters{
		Type:   c.Query("type"),
		Search: c.Query("search"),
		Page:   1,
		Limit:  20,
	}

	// Parse de datas
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filters.StartDate = &startDate
		}
	}
	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filters.EndDate = &endDate
		}
	}

	// Parse de categoria
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			categoryIDUint := uint(categoryID)
			filters.CategoryID = &categoryIDUint
		}
	}

	// Parse de paginação
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			filters.Page = page
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	result, err := services.GetTransactions(userID, filters)
	if err != nil {
		log.Printf("Erro ao buscar transações: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar transações")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transações encontradas", result)
}

// GetTransaction retorna uma transação específica
func GetTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ValidationErrorResponse(c, "ID inválido")
	}

	transaction, err := services.GetTransactionByID(userID, uint(transactionID))
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transação encontrada", transaction)
}

// CreateTransaction cria uma nova transação
func CreateTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Dados inválidos")
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Data inválida (formato esperado: YYYY-MM-DD)")
	}

	var categoryID *uint
	if req.CategoryID != 0 {
		categoryID = &req.CategoryID
	}

	input := services.CreateTransactionInput{
		Type:        req.Type,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        date,
		CategoryID:  categoryID,
		Notes:       req.Notes,
	}

	transaction, err := services.CreateTransaction(userID, input)
	if err != nil {
		log.Printf("Erro ao criar transação: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Transação criada com sucesso", transaction)
}

// UpdateTransaction atualiza uma transação existente
func UpdateTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ValidationErrorResponse(c, "ID inválido")
	}

	var req TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Dados inválidos")
	}

	// Parse da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Data inválida (formato esperado: YYYY-MM-DD)")
	}

	var categoryID *uint
	if req.CategoryID != 0 {
		categoryID = &req.CategoryID
	}

	input := services.CreateTransactionInput{
		Type:        req.Type,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        date,
		CategoryID:  categoryID,
		Notes:       req.Notes,
	}

	transaction, err := services.UpdateTransaction(userID, uint(transactionID), input)
	if err != nil {
		log.Printf("Erro ao atualizar transação: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transação atualizada com sucesso", transaction)
}

// DeleteTransaction exclui uma transação
func DeleteTransaction(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	transactionID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ValidationErrorResponse(c, "ID inválido")
	}

	err = services.DeleteTransaction(userID, uint(transactionID))
	if err != nil {
		log.Printf("Erro ao excluir transação: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transação excluída com sucesso", nil)
}

// ExportTransactions exporta transações para CSV
func ExportTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse dos filtros da query string
	filters := services.TransactionFilters{
		Type:   c.Query("type"),
		Search: c.Query("search"),
	}

	// Parse de datas
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filters.StartDate = &startDate
		}
	}
	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filters.EndDate = &endDate
		}
	}

	// Parse de categoria
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			categoryIDUint := uint(categoryID)
			filters.CategoryID = &categoryIDUint
		}
	}

	csvData, err := services.ExportTransactionsCSV(userID, filters)
	if err != nil {
		log.Printf("Erro ao exportar transações: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao exportar transações")
	}

	// Configurar headers para download de arquivo CSV
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=transacoes.csv")

	return c.Send(csvData)
}
