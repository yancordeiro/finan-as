package services

import (
	"fmt"

	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/models"
)

// GetAllCategories retorna todas as categorias (globais + do usuário)
func GetAllCategories(userID uint) ([]models.Category, error) {
	db := database.GetDB()

	var categories []models.Category

	// Buscar categorias globais (user_id IS NULL) e categorias do usuário
	err := db.Where("user_id IS NULL OR user_id = ?", userID).
		Order("type ASC, name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}

	return categories, nil
}

// GetCategoryByID retorna uma categoria específica
func GetCategoryByID(categoryID uint, userID uint) (*models.Category, error) {
	db := database.GetDB()

	var category models.Category

	// Buscar categoria global ou do usuário
	err := db.Where("id = ? AND (user_id IS NULL OR user_id = ?)", categoryID, userID).
		First(&category).Error

	if err != nil {
		return nil, fmt.Errorf("categoria não encontrada")
	}

	return &category, nil
}

// CreateCategory cria uma nova categoria do usuário
func CreateCategory(userID uint, name, icon, color string, categoryType models.CategoryType) (*models.Category, error) {
	db := database.GetDB()

	// Validações
	if name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	if icon == "" {
		return nil, fmt.Errorf("ícone é obrigatório")
	}
	if color == "" {
		return nil, fmt.Errorf("cor é obrigatória")
	}
	if categoryType != models.CategoryTypeIncome &&
	   categoryType != models.CategoryTypeExpense &&
	   categoryType != models.CategoryTypeBoth {
		return nil, fmt.Errorf("tipo inválido")
	}

	category := models.Category{
		UserID: &userID,
		Name:   name,
		Icon:   icon,
		Color:  color,
		Type:   categoryType,
	}

	if err := db.Create(&category).Error; err != nil {
		return nil, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	return &category, nil
}

// UpdateCategory atualiza uma categoria do usuário
func UpdateCategory(categoryID uint, userID uint, name, icon, color string, categoryType models.CategoryType) (*models.Category, error) {
	db := database.GetDB()

	// Buscar categoria
	var category models.Category
	err := db.Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, fmt.Errorf("categoria não encontrada")
	}

	// Verificar se é uma categoria global (não pode editar)
	if category.IsGlobal() {
		return nil, fmt.Errorf("não é possível editar categorias padrão do sistema")
	}

	// Verificar se pertence ao usuário
	if category.UserID == nil || *category.UserID != userID {
		return nil, fmt.Errorf("você não tem permissão para editar esta categoria")
	}

	// Validações
	if name == "" {
		return nil, fmt.Errorf("nome é obrigatório")
	}
	if icon == "" {
		return nil, fmt.Errorf("ícone é obrigatório")
	}
	if color == "" {
		return nil, fmt.Errorf("cor é obrigatória")
	}

	// Atualizar campos
	category.Name = name
	category.Icon = icon
	category.Color = color
	category.Type = categoryType

	if err := db.Save(&category).Error; err != nil {
		return nil, fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	return &category, nil
}

// DeleteCategory deleta uma categoria do usuário
func DeleteCategory(categoryID uint, userID uint) error {
	db := database.GetDB()

	// Buscar categoria
	var category models.Category
	err := db.Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		return fmt.Errorf("categoria não encontrada")
	}

	// Verificar se é uma categoria global (não pode deletar)
	if category.IsGlobal() {
		return fmt.Errorf("não é possível deletar categorias padrão do sistema")
	}

	// Verificar se pertence ao usuário
	if category.UserID == nil || *category.UserID != userID {
		return fmt.Errorf("você não tem permissão para deletar esta categoria")
	}

	// Verificar se há transações usando esta categoria
	var transactionCount int64
	db.Model(&models.Transaction{}).Where("category_id = ?", categoryID).Count(&transactionCount)

	if transactionCount > 0 {
		return fmt.Errorf("não é possível deletar categoria com transações associadas")
	}

	// Deletar categoria
	if err := db.Delete(&category).Error; err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	return nil
}
