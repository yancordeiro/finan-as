package services

import (
	"fmt"
	"time"

	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/models"
)

// DashboardSummary contém o resumo financeiro do mês
type DashboardSummary struct {
	TotalIncome      int64   `json:"totalIncome"`      // Em centavos
	TotalExpense     int64   `json:"totalExpense"`     // Em centavos
	Balance          int64   `json:"balance"`          // Em centavos
	IncomeChange     float64 `json:"incomeChange"`     // Variação % vs mês anterior
	ExpenseChange    float64 `json:"expenseChange"`    // Variação % vs mês anterior
	TransactionCount int64   `json:"transactionCount"` // Número de transações
	AverageTicket    int64   `json:"averageTicket"`    // Ticket médio em centavos
}

// CategorySummary contém o resumo de uma categoria
type CategorySummary struct {
	CategoryID   uint    `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	CategoryIcon string  `json:"categoryIcon"`
	CategoryColor string `json:"categoryColor"`
	Total        int64   `json:"total"`      // Em centavos
	Count        int64   `json:"count"`      // Número de transações
	Percentage   float64 `json:"percentage"` // Percentual do total
}

// MonthlyFlow contém o fluxo de um mês
type MonthlyFlow struct {
	Month   string `json:"month"`   // Formato: "2024-01"
	Label   string `json:"label"`   // Formato: "Jan/24"
	Income  int64  `json:"income"`  // Em centavos
	Expense int64  `json:"expense"` // Em centavos
	Balance int64  `json:"balance"` // Em centavos
}

// GetDashboardSummary retorna o resumo financeiro do mês
func GetDashboardSummary(userID uint, month int, year int) (*DashboardSummary, error) {
	db := database.GetDB()

	// Definir período do mês atual
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// Buscar totais do mês atual
	var currentIncome, currentExpense int64
	var transactionCount int64

	// Total de entradas
	db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeIncome, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&currentIncome)

	// Total de saídas
	db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeExpense, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&currentExpense)

	// Contagem de transações
	db.Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date < ?",
			userID, startDate, endDate).
		Count(&transactionCount)

	// Buscar totais do mês anterior para calcular variação
	prevStartDate := startDate.AddDate(0, -1, 0)
	prevEndDate := startDate

	var prevIncome, prevExpense int64

	db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeIncome, prevStartDate, prevEndDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&prevIncome)

	db.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeExpense, prevStartDate, prevEndDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&prevExpense)

	// Calcular variação percentual
	var incomeChange, expenseChange float64

	if prevIncome > 0 {
		incomeChange = float64(currentIncome-prevIncome) / float64(prevIncome) * 100
	} else if currentIncome > 0 {
		incomeChange = 100 // Se não tinha antes e agora tem, é 100% de aumento
	}

	if prevExpense > 0 {
		expenseChange = float64(currentExpense-prevExpense) / float64(prevExpense) * 100
	} else if currentExpense > 0 {
		expenseChange = 100
	}

	// Calcular ticket médio
	var averageTicket int64
	if transactionCount > 0 {
		averageTicket = (currentIncome + currentExpense) / transactionCount
	}

	return &DashboardSummary{
		TotalIncome:      currentIncome,
		TotalExpense:     currentExpense,
		Balance:          currentIncome - currentExpense,
		IncomeChange:     incomeChange,
		ExpenseChange:    expenseChange,
		TransactionCount: transactionCount,
		AverageTicket:    averageTicket,
	}, nil
}

// GetExpensesByCategory retorna gastos agrupados por categoria
func GetExpensesByCategory(userID uint, month int, year int) ([]CategorySummary, error) {
	db := database.GetDB()

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	type categoryResult struct {
		CategoryID uint
		Total      int64
		Count      int64
	}

	var results []categoryResult

	err := db.Model(&models.Transaction{}).
		Select("category_id, SUM(amount) as total, COUNT(*) as count").
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeExpense, startDate, endDate).
		Group("category_id").
		Order("total DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar gastos por categoria: %w", err)
	}

	// Calcular total geral para percentual
	var totalExpense int64
	for _, r := range results {
		totalExpense += r.Total
	}

	// Buscar informações das categorias e calcular percentual
	summaries := make([]CategorySummary, 0, len(results))
	for _, r := range results {
		var category models.Category
		db.First(&category, r.CategoryID)

		percentage := float64(0)
		if totalExpense > 0 {
			percentage = float64(r.Total) / float64(totalExpense) * 100
		}

		summaries = append(summaries, CategorySummary{
			CategoryID:    r.CategoryID,
			CategoryName:  category.Name,
			CategoryIcon:  category.Icon,
			CategoryColor: category.Color,
			Total:         r.Total,
			Count:         r.Count,
			Percentage:    percentage,
		})
	}

	return summaries, nil
}

// GetIncomesByCategory retorna entradas agrupadas por categoria
func GetIncomesByCategory(userID uint, month int, year int) ([]CategorySummary, error) {
	db := database.GetDB()

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	type categoryResult struct {
		CategoryID uint
		Total      int64
		Count      int64
	}

	var results []categoryResult

	err := db.Model(&models.Transaction{}).
		Select("category_id, SUM(amount) as total, COUNT(*) as count").
		Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, models.TransactionTypeIncome, startDate, endDate).
		Group("category_id").
		Order("total DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar entradas por categoria: %w", err)
	}

	// Calcular total geral para percentual
	var totalIncome int64
	for _, r := range results {
		totalIncome += r.Total
	}

	// Buscar informações das categorias e calcular percentual
	summaries := make([]CategorySummary, 0, len(results))
	for _, r := range results {
		var category models.Category
		db.First(&category, r.CategoryID)

		percentage := float64(0)
		if totalIncome > 0 {
			percentage = float64(r.Total) / float64(totalIncome) * 100
		}

		summaries = append(summaries, CategorySummary{
			CategoryID:    r.CategoryID,
			CategoryName:  category.Name,
			CategoryIcon:  category.Icon,
			CategoryColor: category.Color,
			Total:         r.Total,
			Count:         r.Count,
			Percentage:    percentage,
		})
	}

	return summaries, nil
}

// GetMonthlyFlow retorna o fluxo financeiro dos últimos N meses
func GetMonthlyFlow(userID uint, months int) ([]MonthlyFlow, error) {
	db := database.GetDB()

	if months < 1 {
		months = 6
	}
	if months > 12 {
		months = 12
	}

	now := time.Now()
	results := make([]MonthlyFlow, 0, months)

	// Nomes dos meses em português
	monthNames := []string{
		"", "Jan", "Fev", "Mar", "Abr", "Mai", "Jun",
		"Jul", "Ago", "Set", "Out", "Nov", "Dez",
	}

	for i := months - 1; i >= 0; i-- {
		// Calcular o mês
		targetDate := now.AddDate(0, -i, 0)
		year := targetDate.Year()
		month := int(targetDate.Month())

		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		var income, expense int64

		db.Model(&models.Transaction{}).
			Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
				userID, models.TransactionTypeIncome, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&income)

		db.Model(&models.Transaction{}).
			Where("user_id = ? AND type = ? AND date >= ? AND date < ?",
				userID, models.TransactionTypeExpense, startDate, endDate).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&expense)

		results = append(results, MonthlyFlow{
			Month:   fmt.Sprintf("%d-%02d", year, month),
			Label:   fmt.Sprintf("%s/%02d", monthNames[month], year%100),
			Income:  income,
			Expense: expense,
			Balance: income - expense,
		})
	}

	return results, nil
}

// GetRecentTransactions retorna as últimas transações
func GetRecentTransactions(userID uint, limit int) ([]models.Transaction, error) {
	db := database.GetDB()

	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var transactions []models.Transaction

	err := db.Preload("Category").
		Where("user_id = ?", userID).
		Order("date DESC, created_at DESC").
		Limit(limit).
		Find(&transactions).Error

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar transações recentes: %w", err)
	}

	return transactions, nil
}
