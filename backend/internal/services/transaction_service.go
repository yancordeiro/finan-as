package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/models"
)

// TransactionFilters contém os filtros para busca de transações
type TransactionFilters struct {
	StartDate  *time.Time
	EndDate    *time.Time
	Type       string
	CategoryID *uint
	Search     string
	Page       int
	Limit      int
}

// TransactionListResult contém o resultado da listagem com paginação
type TransactionListResult struct {
	Transactions []models.Transaction `json:"transactions"`
	Total        int64                `json:"total"`
	Page         int                  `json:"page"`
	Limit        int                  `json:"limit"`
	TotalPages   int                  `json:"totalPages"`
}

// CreateTransactionInput contém os dados para criar uma transação
type CreateTransactionInput struct {
	Type        models.TransactionType `json:"type"`
	Description string                 `json:"description"`
	Amount      int64                  `json:"amount"` // Em centavos
	Date        time.Time              `json:"date"`
	CategoryID  *uint                  `json:"categoryId"` // Opcional - permite transações sem categoria
	Notes       string                 `json:"notes"`
}

// GetTransactions retorna transações do usuário com filtros e paginação
func GetTransactions(userID uint, filters TransactionFilters) (*TransactionListResult, error) {
	db := database.GetDB()

	// Valores padrão para paginação
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	// Query base
	query := db.Model(&models.Transaction{}).Where("user_id = ?", userID)

	// Aplicar filtros
	if filters.StartDate != nil {
		query = query.Where("date >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		// Adicionar 1 dia para incluir o dia final completo
		endDate := filters.EndDate.Add(24 * time.Hour)
		query = query.Where("date < ?", endDate)
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.Search != "" {
		query = query.Where("description ILIKE ?", "%"+filters.Search+"%")
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("erro ao contar transações: %w", err)
	}

	// Calcular total de páginas
	totalPages := int(total) / filters.Limit
	if int(total)%filters.Limit > 0 {
		totalPages++
	}

	// Buscar transações com paginação
	var transactions []models.Transaction
	offset := (filters.Page - 1) * filters.Limit

	err := query.
		Preload("Category").
		Order("date DESC, created_at DESC").
		Offset(offset).
		Limit(filters.Limit).
		Find(&transactions).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transações: %w", err)
	}

	return &TransactionListResult{
		Transactions: transactions,
		Total:        total,
		Page:         filters.Page,
		Limit:        filters.Limit,
		TotalPages:   totalPages,
	}, nil
}

// GetTransactionByID retorna uma transação específica do usuário
func GetTransactionByID(userID uint, transactionID uint) (*models.Transaction, error) {
	db := database.GetDB()

	var transaction models.Transaction
	err := db.Preload("Category").
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error

	if err != nil {
		return nil, fmt.Errorf("transação não encontrada")
	}

	return &transaction, nil
}

// CreateTransaction cria uma nova transação
func CreateTransaction(userID uint, input CreateTransactionInput) (*models.Transaction, error) {
	db := database.GetDB()

	// Validações
	if input.Description == "" {
		return nil, fmt.Errorf("descrição é obrigatória")
	}
	if input.Amount <= 0 {
		return nil, fmt.Errorf("valor deve ser maior que zero")
	}
	if input.Type != models.TransactionTypeIncome && input.Type != models.TransactionTypeExpense {
		return nil, fmt.Errorf("tipo inválido (deve ser 'income' ou 'expense')")
	}

	// Verificar categoria apenas se fornecida
	if input.CategoryID != nil {
		category, err := GetCategoryByID(*input.CategoryID, userID)
		if err != nil {
			return nil, fmt.Errorf("categoria não encontrada")
		}

		// Verificar se a categoria é compatível com o tipo de transação
		if category.Type != models.CategoryTypeBoth {
			if input.Type == models.TransactionTypeIncome && category.Type != models.CategoryTypeIncome {
				return nil, fmt.Errorf("categoria não permite transações de entrada")
			}
			if input.Type == models.TransactionTypeExpense && category.Type != models.CategoryTypeExpense {
				return nil, fmt.Errorf("categoria não permite transações de saída")
			}
		}
	}

	transaction := models.Transaction{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Type:        input.Type,
		Description: input.Description,
		Amount:      input.Amount,
		Date:        input.Date,
		Notes:       input.Notes,
	}

	if err := db.Create(&transaction).Error; err != nil {
		return nil, fmt.Errorf("erro ao criar transação: %w", err)
	}

	// Carregar categoria para retornar
	db.Preload("Category").First(&transaction, transaction.ID)

	return &transaction, nil
}

// UpdateTransaction atualiza uma transação existente
func UpdateTransaction(userID uint, transactionID uint, input CreateTransactionInput) (*models.Transaction, error) {
	db := database.GetDB()

	// Buscar transação
	var transaction models.Transaction
	err := db.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error
	if err != nil {
		return nil, fmt.Errorf("transação não encontrada")
	}

	// Validações
	if input.Description == "" {
		return nil, fmt.Errorf("descrição é obrigatória")
	}
	if input.Amount <= 0 {
		return nil, fmt.Errorf("valor deve ser maior que zero")
	}
	if input.Type != models.TransactionTypeIncome && input.Type != models.TransactionTypeExpense {
		return nil, fmt.Errorf("tipo inválido")
	}

	// Verificar categoria apenas se fornecida
	if input.CategoryID != nil {
		category, err := GetCategoryByID(*input.CategoryID, userID)
		if err != nil {
			return nil, fmt.Errorf("categoria não encontrada")
		}

		// Verificar compatibilidade da categoria
		if category.Type != models.CategoryTypeBoth {
			if input.Type == models.TransactionTypeIncome && category.Type != models.CategoryTypeIncome {
				return nil, fmt.Errorf("categoria não permite transações de entrada")
			}
			if input.Type == models.TransactionTypeExpense && category.Type != models.CategoryTypeExpense {
				return nil, fmt.Errorf("categoria não permite transações de saída")
			}
		}
	}

	// Atualizar campos
	transaction.Type = input.Type
	transaction.Description = input.Description
	transaction.Amount = input.Amount
	transaction.Date = input.Date
	transaction.CategoryID = input.CategoryID
	transaction.Notes = input.Notes

	if err := db.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("erro ao atualizar transação: %w", err)
	}

	// Carregar categoria para retornar
	db.Preload("Category").First(&transaction, transaction.ID)

	return &transaction, nil
}

// DeleteTransaction exclui uma transação
func DeleteTransaction(userID uint, transactionID uint) error {
	db := database.GetDB()

	result := db.Where("id = ? AND user_id = ?", transactionID, userID).Delete(&models.Transaction{})

	if result.Error != nil {
		return fmt.Errorf("erro ao excluir transação: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transação não encontrada")
	}

	return nil
}

// ExportTransactionsCSV exporta transações para CSV
func ExportTransactionsCSV(userID uint, filters TransactionFilters) ([]byte, error) {
	// Remover paginação para exportar tudo
	filters.Page = 1
	filters.Limit = 10000 // Limite alto para exportação

	result, err := GetTransactions(userID, filters)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Cabeçalho
	header := []string{"Data", "Tipo", "Descrição", "Categoria", "Valor (R$)", "Observações"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("erro ao escrever cabeçalho CSV: %w", err)
	}

	// Dados
	for _, t := range result.Transactions {
		tipoStr := "Entrada"
		if t.Type == models.TransactionTypeExpense {
			tipoStr = "Saída"
		}

		categoryName := ""
		if t.Category.ID != 0 {
			categoryName = t.Category.Name
		}

		row := []string{
			t.Date.Format("02/01/2006"),
			tipoStr,
			t.Description,
			categoryName,
			fmt.Sprintf("%.2f", t.AmountInReais()),
			t.Notes,
		}

		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("erro ao escrever linha CSV: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("erro ao finalizar CSV: %w", err)
	}

	return buf.Bytes(), nil
}
