package handlers

import (
	"bytes"
	"log"
	"strconv"

	"github.com/financas/backend/internal/middleware"
	"github.com/financas/backend/internal/services"
	"github.com/financas/backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// ImportOFX importa transações de um arquivo OFX
func ImportOFX(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Obter arquivo do form
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Arquivo OFX não enviado")
	}

	// Verificar extensão
	if len(file.Filename) < 4 || (file.Filename[len(file.Filename)-4:] != ".ofx" && file.Filename[len(file.Filename)-4:] != ".OFX") {
		return utils.ValidationErrorResponse(c, "Formato de arquivo inválido. Envie um arquivo .ofx")
	}

	// Obter categoria padrão (opcional)
	var categoryID *uint
	categoryIDStr := c.FormValue("categoryId")
	if categoryIDStr != "" {
		parsedID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			return utils.ValidationErrorResponse(c, "ID de categoria inválido")
		}
		id := uint(parsedID)
		categoryID = &id
	}

	// Abrir arquivo
	f, err := file.Open()
	if err != nil {
		log.Printf("Erro ao abrir arquivo OFX: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao processar arquivo")
	}
	defer f.Close()

	// Ler conteúdo do arquivo
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		log.Printf("Erro ao ler arquivo OFX: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao ler arquivo")
	}

	// Importar transações
	result, err := services.ImportOFX(userID, bytes.NewReader(buf.Bytes()), categoryID)
	if err != nil {
		log.Printf("Erro ao importar OFX: %v", err)
		return utils.ValidationErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Importação concluída", result)
}

// PreviewOFX faz preview das transações do arquivo OFX sem importar
func PreviewOFX(c *fiber.Ctx) error {
	// Obter arquivo do form
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Arquivo OFX não enviado")
	}

	// Verificar extensão
	if len(file.Filename) < 4 || (file.Filename[len(file.Filename)-4:] != ".ofx" && file.Filename[len(file.Filename)-4:] != ".OFX") {
		return utils.ValidationErrorResponse(c, "Formato de arquivo inválido. Envie um arquivo .ofx")
	}

	// Abrir arquivo
	f, err := file.Open()
	if err != nil {
		log.Printf("Erro ao abrir arquivo OFX: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao processar arquivo")
	}
	defer f.Close()

	// Ler conteúdo do arquivo
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		log.Printf("Erro ao ler arquivo OFX: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao ler arquivo")
	}

	// Parse do arquivo
	transactions, err := services.ParseOFX(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Printf("Erro ao fazer parse do OFX: %v", err)
		return utils.ValidationErrorResponse(c, "Erro ao processar arquivo OFX")
	}

	// Formatar resposta
	type PreviewTransaction struct {
		FITID       string `json:"fitid"`
		Date        string `json:"date"`
		Description string `json:"description"`
		Amount      int64  `json:"amount"`
		Type        string `json:"type"`
	}

	preview := make([]PreviewTransaction, 0, len(transactions))
	for _, t := range transactions {
		transType := "income"
		amount := t.Amount
		if amount < 0 {
			transType = "expense"
			amount = -amount
		}

		description := t.Name
		if description == "" {
			description = t.Memo
		}

		preview = append(preview, PreviewTransaction{
			FITID:       t.FITID,
			Date:        t.Date.Format("2006-01-02"),
			Description: description,
			Amount:      amount,
			Type:        transType,
		})
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Preview do arquivo OFX", fiber.Map{
		"transactions": preview,
		"total":        len(preview),
	})
}
