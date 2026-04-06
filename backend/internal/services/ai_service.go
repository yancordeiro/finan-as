package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/financas/backend/internal/models"
	"gorm.io/gorm"
)

type AIService struct {
	db          *gorm.DB
	ollamaURL   string
	ollamaModel string
	httpClient  *http.Client
}

func NewAIService(db *gorm.DB) *AIService {
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434" // Default para Docker
	}

	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "llama3.1:8b" // Default
	}

	return &AIService{
		db:          db,
		ollamaURL:   ollamaURL,
		ollamaModel: ollamaModel,
		httpClient: &http.Client{
			Timeout: 180 * time.Second, // Ollama pode demorar em CPU (3 minutos)
		},
	}
}

// Chat processa mensagem do usuário e retorna resposta da IA
func (s *AIService) Chat(userID uint, message string) (*models.AIResponse, error) {
	// 1. Buscar contexto financeiro do usuário
	context, err := s.buildFinancialContext(userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contexto financeiro: %w", err)
	}

	// 2. Montar prompt com contexto
	systemPrompt := s.buildSystemPrompt(context)

	// 3. Chamar Ollama
	aiResponse, err := s.callOllama(systemPrompt, message)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar Ollama: %w", err)
	}

	// 4. Montar resposta
	response := &models.AIResponse{
		Response:         aiResponse,
		FinancialContext: context,
		Suggestions:      s.extractSuggestions(aiResponse, context),
	}

	return response, nil
}

// buildFinancialContext busca dados do usuário e monta contexto
func (s *AIService) buildFinancialContext(userID uint) (*models.FinancialContext, error) {
	// Buscar transações dos últimos 30 dias
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var transactions []models.Transaction
	err := s.db.Preload("Category").
		Where("user_id = ? AND date >= ?", userID, thirtyDaysAgo).
		Order("date DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	// Buscar categorias do usuário
	var categories []models.Category
	err = s.db.Where("user_id = ?", userID).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// Calcular estatísticas
	var totalIncome, totalExpense int64
	categoryTotals := make(map[string]*models.CategoryStats)

	for _, t := range transactions {
		amount := t.Amount

		if t.Type == models.TransactionTypeIncome {
			totalIncome += amount
		} else {
			totalExpense += amount
		}

		// Agrupar por categoria
		if t.Category != nil {
			catName := t.Category.Name
			if _, exists := categoryTotals[catName]; !exists {
				categoryTotals[catName] = &models.CategoryStats{
					Amount: 0,
					Count:  0,
				}
			}
			categoryTotals[catName].Amount += float64(amount) / 100.0
			categoryTotals[catName].Count++
		}
	}

	// Calcular percentuais
	totalExpenseFloat := float64(totalExpense) / 100.0
	for _, stats := range categoryTotals {
		if totalExpenseFloat > 0 {
			stats.Percentage = (stats.Amount / totalExpenseFloat) * 100
		}
	}

	context := &models.FinancialContext{
		TotalIncome:       float64(totalIncome) / 100.0,
		TotalExpense:      float64(totalExpense) / 100.0,
		Balance:           float64(totalIncome-totalExpense) / 100.0,
		CategoryBreakdown: categoryTotals,
		Period:            "last_30_days",
		TransactionCount:  len(transactions),
	}

	return context, nil
}

// buildSystemPrompt monta prompt do sistema com contexto financeiro
func (s *AIService) buildSystemPrompt(context *models.FinancialContext) string {
	prompt := `Você é um assistente financeiro pessoal especializado e experiente.

CONTEXTO FINANCEIRO DO USUÁRIO (últimos 30 dias):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Resumo Geral:
• Receitas: R$ %.2f
• Despesas: R$ %.2f
• Saldo: R$ %.2f
• Total de transações: %d

Gastos por Categoria:
`
	prompt = fmt.Sprintf(prompt,
		context.TotalIncome,
		context.TotalExpense,
		context.Balance,
		context.TransactionCount,
	)

	// Adicionar breakdown por categoria
	for catName, stats := range context.CategoryBreakdown {
		prompt += fmt.Sprintf("• %s: R$ %.2f (%.1f%% do total) - %d transações\n",
			catName,
			stats.Amount,
			stats.Percentage,
			stats.Count,
		)
	}

	prompt += `
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SUAS CAPACIDADES:
✓ Analisar padrões de gastos baseado nos dados reais acima
✓ Identificar oportunidades de economia
✓ Sugerir metas realistas de economia
✓ Detectar gastos anormais ou excessivos
✓ Criar planos de ação personalizados
✓ Responder perguntas sobre finanças pessoais

INSTRUÇÕES IMPORTANTES:
• Use SEMPRE os dados reais fornecidos acima
• Seja específico com números e valores reais
• Dê recomendações práticas e acionáveis
• Use uma linguagem amigável mas profissional
• Seja direto e objetivo
• Evite respostas genéricas
• Foque em ações concretas que o usuário pode tomar

FORMATO DE RESPOSTA:
• Use markdown para formatação
• Destaque valores monetários com **negrito**
• Use bullet points para listas
• Seja conciso (máximo 200 palavras por resposta)

Responda à pergunta do usuário baseando-se nos dados financeiros reais dele:
`

	return prompt
}

// callOllama faz chamada HTTP para o Ollama
func (s *AIService) callOllama(systemPrompt, userMessage string) (string, error) {
	messages := []models.OllamaMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	reqBody := models.OllamaChatRequest{
		Model:    s.ollamaModel,
		Messages: messages,
		Stream:   false,
		Options: &models.OllamaOptions{
			Temperature: 0.7,
			NumPredict:  500, // Limitar para respostas concisas
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %w", err)
	}

	// Fazer request para Ollama
	url := fmt.Sprintf("%s/api/chat", s.ollamaURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao conectar com Ollama em %s: %w (Verifique se Ollama está rodando)", s.ollamaURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama retornou erro %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp models.OllamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// extractSuggestions extrai sugestões acionáveis da resposta (básico por enquanto)
func (s *AIService) extractSuggestions(response string, context *models.FinancialContext) []models.AISuggestion {
	suggestions := []models.AISuggestion{}

	// Sugestão simples: se saldo negativo, sugerir reduzir categoria mais alta
	if context.Balance < 0 {
		// Encontrar categoria com maior gasto
		var highestCat string
		var highestAmount float64

		for catName, stats := range context.CategoryBreakdown {
			if stats.Amount > highestAmount {
				highestAmount = stats.Amount
				highestCat = catName
			}
		}

		if highestCat != "" {
			suggestions = append(suggestions, models.AISuggestion{
				Type:        "reduce_category",
				Title:       fmt.Sprintf("Reduzir gastos em %s", highestCat),
				Description: fmt.Sprintf("Você gastou R$ %.2f em %s. Reduzir 20%% economizaria R$ %.2f/mês", highestAmount, highestCat, highestAmount*0.2),
				Action:      "/transactions?category=" + highestCat,
			})
		}
	}

	// Mais sugestões podem ser adicionadas aqui no futuro

	return suggestions
}

// GetFinancialSummary retorna apenas o contexto financeiro (útil para dashboard)
func (s *AIService) GetFinancialSummary(userID uint) (*models.FinancialContext, error) {
	return s.buildFinancialContext(userID)
}
