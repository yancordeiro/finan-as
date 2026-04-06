package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse retorna uma resposta de sucesso padronizada
func SuccessResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse retorna uma resposta de erro padronizada
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error:   message,
	})
}

// ValidationErrorResponse retorna uma resposta de erro de validação
func ValidationErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// UnauthorizedResponse retorna uma resposta de não autorizado
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Não autorizado"
	}
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// NotFoundResponse retorna uma resposta de não encontrado
func NotFoundResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Recurso não encontrado"
	}
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// InternalErrorResponse retorna uma resposta de erro interno
func InternalErrorResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Erro interno do servidor"
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}
