package models

// AIRequest representa uma mensagem do usuário para a IA
type AIRequest struct {
	Message string `json:"message" validate:"required,min=1,max=1000"`
}

// AIResponse representa a resposta da IA
type AIResponse struct {
	Response         string              `json:"response"`
	Suggestions      []AISuggestion      `json:"suggestions,omitempty"`
	FinancialContext *FinancialContext   `json:"financialContext,omitempty"`
}

// AISuggestion representa uma ação sugerida pela IA
type AISuggestion struct {
	Type        string  `json:"type"` // "create_goal", "reduce_category", "review_transactions"
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Action      string  `json:"action"` // URL ou comando frontend deve executar
	Data        any     `json:"data,omitempty"`
}

// FinancialContext representa o contexto financeiro do usuário usado pela IA
type FinancialContext struct {
	TotalIncome       float64                     `json:"totalIncome"`
	TotalExpense      float64                     `json:"totalExpense"`
	Balance           float64                     `json:"balance"`
	CategoryBreakdown map[string]*CategoryStats   `json:"categoryBreakdown"`
	Period            string                      `json:"period"` // "last_30_days", etc
	TransactionCount  int                         `json:"transactionCount"`
}

// CategoryStats estatísticas por categoria
type CategoryStats struct {
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
	Count      int     `json:"count"`
}

// OllamaChatRequest representa request para API do Ollama
type OllamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []OllamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  *OllamaOptions  `json:"options,omitempty"`
}

// OllamaMessage representa uma mensagem no formato Ollama
type OllamaMessage struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"`
}

// OllamaOptions opções para geração de texto
type OllamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"` // 0.0 a 1.0
	NumPredict  int     `json:"num_predict,omitempty"` // Max tokens
}

// OllamaChatResponse resposta da API Ollama
type OllamaChatResponse struct {
	Model     string        `json:"model"`
	CreatedAt string        `json:"created_at"`
	Message   OllamaMessage `json:"message"`
	Done      bool          `json:"done"`
}
