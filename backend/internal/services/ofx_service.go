package services

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/models"
)

// OFXTransaction representa uma transação extraída do arquivo OFX
type OFXTransaction struct {
	FITID    string    // ID único da transação no banco
	Type     string    // DEBIT ou CREDIT
	Date     time.Time // Data da transação
	Amount   int64     // Valor em centavos
	Name     string    // Nome/descrição
	Memo     string    // Memo adicional
}

// OFXImportResult representa o resultado da importação
type OFXImportResult struct {
	TotalRead    int `json:"totalRead"`
	Imported     int `json:"imported"`
	Skipped      int `json:"skipped"`
	Duplicates   int `json:"duplicates"`
}

// ParseOFX faz o parsing de um arquivo OFX e retorna as transações
func ParseOFX(reader io.Reader) ([]OFXTransaction, error) {
	var transactions []OFXTransaction

	scanner := bufio.NewScanner(reader)
	var content strings.Builder

	// Ler todo o conteúdo
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo OFX: %w", err)
	}

	text := content.String()

	// Encontrar todas as transações (STMTTRN)
	stmtTrnRegex := regexp.MustCompile(`(?s)<STMTTRN>(.*?)</STMTTRN>`)
	matches := stmtTrnRegex.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		trn := match[1]

		transaction := OFXTransaction{}

		// Extrair TRNTYPE (DEBIT, CREDIT, etc)
		if val := extractOFXValue(trn, "TRNTYPE"); val != "" {
			transaction.Type = val
		}

		// Extrair DTPOSTED (data)
		if val := extractOFXValue(trn, "DTPOSTED"); val != "" {
			if date, err := parseOFXDate(val); err == nil {
				transaction.Date = date
			}
		}

		// Extrair TRNAMT (valor)
		if val := extractOFXValue(trn, "TRNAMT"); val != "" {
			if amount, err := parseOFXAmount(val); err == nil {
				transaction.Amount = amount
			}
		}

		// Extrair FITID (ID único)
		if val := extractOFXValue(trn, "FITID"); val != "" {
			transaction.FITID = val
		}

		// Extrair NAME (descrição)
		if val := extractOFXValue(trn, "NAME"); val != "" {
			transaction.Name = strings.TrimSpace(val)
		}

		// Extrair MEMO (memo adicional)
		if val := extractOFXValue(trn, "MEMO"); val != "" {
			transaction.Memo = strings.TrimSpace(val)
		}

		// Só adicionar se tiver os campos essenciais
		if transaction.FITID != "" && transaction.Amount != 0 {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

// extractOFXValue extrai o valor de uma tag OFX
func extractOFXValue(content, tag string) string {
	// Padrão 1: <TAG>valor</TAG>
	regex1 := regexp.MustCompile(`<` + tag + `>([^<]*)<\/` + tag + `>`)
	if match := regex1.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	// Padrão 2: <TAG>valor (sem fechamento, comum em OFX v1)
	regex2 := regexp.MustCompile(`<` + tag + `>([^\r\n<]+)`)
	if match := regex2.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return ""
}

// parseOFXDate converte data OFX (YYYYMMDDHHMMSS) para time.Time
func parseOFXDate(dateStr string) (time.Time, error) {
	// Remover timezone se presente (ex: [-3:GMT])
	if idx := strings.Index(dateStr, "["); idx > 0 {
		dateStr = dateStr[:idx]
	}

	// Tentar diferentes formatos
	formats := []string{
		"20060102150405",
		"20060102120000",
		"20060102",
	}

	for _, format := range formats {
		if len(dateStr) >= len(format) {
			if t, err := time.Parse(format, dateStr[:len(format)]); err == nil {
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("formato de data inválido: %s", dateStr)
}

// parseOFXAmount converte valor OFX para centavos
func parseOFXAmount(amountStr string) (int64, error) {
	// Limpar o valor (remover espaços)
	amountStr = strings.TrimSpace(amountStr)

	// Substituir vírgula por ponto (formato brasileiro)
	amountStr = strings.ReplaceAll(amountStr, ",", ".")

	// Converter para float
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("valor inválido: %s", amountStr)
	}

	// Converter para centavos
	return int64(amount * 100), nil
}

// ImportOFX importa transações de um arquivo OFX
func ImportOFX(userID uint, reader io.Reader, defaultCategoryID *uint) (*OFXImportResult, error) {
	db := database.GetDB()

	// Parse do arquivo
	ofxTransactions, err := ParseOFX(reader)
	if err != nil {
		return nil, err
	}

	result := &OFXImportResult{
		TotalRead: len(ofxTransactions),
	}

	// Verificar se a categoria padrão existe (se fornecida)
	if defaultCategoryID != nil {
		_, err = GetCategoryByID(*defaultCategoryID, userID)
		if err != nil {
			return nil, fmt.Errorf("categoria padrão não encontrada")
		}
	}

	for _, ofxTrn := range ofxTransactions {
		// Verificar se já existe (deduplicação por FITID)
		var existingCount int64
		db.Model(&models.Transaction{}).
			Where("user_id = ? AND ofx_id = ?", userID, ofxTrn.FITID).
			Count(&existingCount)

		if existingCount > 0 {
			result.Duplicates++
			continue
		}

		// Determinar tipo da transação
		var transactionType models.TransactionType
		amount := ofxTrn.Amount

		if amount < 0 {
			transactionType = models.TransactionTypeExpense
			amount = -amount // Converter para positivo
		} else {
			transactionType = models.TransactionTypeIncome
		}

		// Também verificar pelo TRNTYPE
		if strings.ToUpper(ofxTrn.Type) == "DEBIT" {
			transactionType = models.TransactionTypeExpense
			if amount < 0 {
				amount = -amount
			}
		} else if strings.ToUpper(ofxTrn.Type) == "CREDIT" {
			transactionType = models.TransactionTypeIncome
			if amount < 0 {
				amount = -amount
			}
		}

		// Determinar descrição
		description := ofxTrn.Name
		if description == "" {
			description = ofxTrn.Memo
		}
		if description == "" {
			description = "Transação OFX"
		}

		// Criar transação
		transaction := models.Transaction{
			UserID:      userID,
			CategoryID:  defaultCategoryID,
			Type:        transactionType,
			Description: description,
			Amount:      amount,
			Date:        ofxTrn.Date,
			Notes:       ofxTrn.Memo,
			OFXID:       ofxTrn.FITID,
		}

		if err := db.Create(&transaction).Error; err != nil {
			result.Skipped++
			continue
		}

		result.Imported++
	}

	return result, nil
}
