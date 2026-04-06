package handlers

import (
	"log"
	"strconv"

	"github.com/financas/backend/internal/middleware"
	"github.com/financas/backend/internal/models"
	"github.com/financas/backend/internal/services"
	"github.com/financas/backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// CategoryRequest representa o corpo da requisição para criar/atualizar categoria
type CategoryRequest struct {
	Name string                `json:"name"`
	Icon string                `json:"icon"`
	Color string               `json:"color"`
	Type models.CategoryType   `json:"type"`
}

// GetCategories retorna todas as categorias (globais + do usuário)
func GetCategories(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	categories, err := services.GetAllCategories(userID)
	if err != nil {
		log.Printf("Erro ao buscar categorias: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao buscar categorias")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Categorias encontradas", categories)
}

// GetCategory retorna uma categoria específica
func GetCategory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse do ID
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ValidationErrorResponse(c, "ID inválido")
	}

	category, err := services.GetCategoryByID(uint(categoryID), userID)
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Categoria encontrada", category)
}

// CreateCategory cria uma nova categoria
func CreateCategory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Dados inválidos")
	}

	category, err := services.CreateCategory(userID, req.Name, req.Icon, req.Color, req.Type)
	if err != nil {
		log.Printf("Erro ao criar categoria: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Categoria criada com sucesso", category)
}

// UpdateCategory atualiza uma categoria existente
func UpdateCategory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse do ID
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ValidationErrorResponse(c, "ID inválido")
	}

	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Dados inválidos")
	}

	category, err := services.UpdateCategory(uint(categoryID), userID, req.Name, req.Icon, req.Color, req.Type)
	if err != nil {
		log.Printf("Erro ao atualizar categoria: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Categoria atualizada com sucesso", category)
}

// DeleteCategory deleta uma categoria
func DeleteCategory(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Parse do ID
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ValidationErrorResponse(c, "ID inválido")
	}

	err = services.DeleteCategory(uint(categoryID), userID)
	if err != nil {
		log.Printf("Erro ao deletar categoria: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Categoria deletada com sucesso", nil)
}
