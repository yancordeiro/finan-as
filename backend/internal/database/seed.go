package database

import (
	"log"

	"github.com/financas/backend/internal/models"
)

// SeedDefaultCategories cria as categorias padrão do sistema
func SeedDefaultCategories() error {
	log.Println("Verificando categorias padrão...")

	// Verificar se já existem categorias globais
	var count int64
	DB.Model(&models.Category{}).Where("user_id IS NULL").Count(&count)

	if count > 0 {
		log.Printf("✓ Categorias padrão já existem (%d categorias)", count)
		return nil
	}

	log.Println("Criando categorias padrão...")

	// Categorias de despesas
	expenseCategories := []models.Category{
		{
			UserID: nil,
			Name:   "Alimentação",
			Icon:   "🍔",
			Color:  "#FF6B6B",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Transporte",
			Icon:   "🚗",
			Color:  "#4ECDC4",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Saúde",
			Icon:   "🏥",
			Color:  "#45B7D1",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Lazer",
			Icon:   "🎮",
			Color:  "#A259FF",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Moradia",
			Icon:   "🏠",
			Color:  "#FFA07A",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Educação",
			Icon:   "📚",
			Color:  "#95E1D3",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Compras",
			Icon:   "🛒",
			Color:  "#FFB6C1",
			Type:   models.CategoryTypeExpense,
		},
		{
			UserID: nil,
			Name:   "Contas",
			Icon:   "📄",
			Color:  "#FFD700",
			Type:   models.CategoryTypeExpense,
		},
	}

	// Categorias de receitas
	incomeCategories := []models.Category{
		{
			UserID: nil,
			Name:   "Salário",
			Icon:   "💼",
			Color:  "#00E676",
			Type:   models.CategoryTypeIncome,
		},
		{
			UserID: nil,
			Name:   "Freelance",
			Icon:   "💻",
			Color:  "#00D4FF",
			Type:   models.CategoryTypeIncome,
		},
		{
			UserID: nil,
			Name:   "Investimentos",
			Icon:   "📈",
			Color:  "#FFD93D",
			Type:   models.CategoryTypeIncome,
		},
		{
			UserID: nil,
			Name:   "Bônus",
			Icon:   "🎁",
			Color:  "#6BCF7F",
			Type:   models.CategoryTypeIncome,
		},
	}

	// Categoria genérica
	otherCategory := models.Category{
		UserID: nil,
		Name:   "Outros",
		Icon:   "📦",
		Color:  "#7A7A9A",
		Type:   models.CategoryTypeBoth,
	}

	// Inserir todas as categorias
	allCategories := append(expenseCategories, incomeCategories...)
	allCategories = append(allCategories, otherCategory)

	for _, category := range allCategories {
		if err := DB.Create(&category).Error; err != nil {
			return err
		}
	}

	log.Printf("✓ %d categorias padrão criadas com sucesso", len(allCategories))
	return nil
}
