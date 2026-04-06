package middleware

import (
	"strings"

	"github.com/financas/backend/internal/config"
	"github.com/financas/backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware valida o JWT e injeta o userID no contexto
func AuthMiddleware(c *fiber.Ctx) error {
	// Extrair o token do header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.UnauthorizedResponse(c, "Token de autenticação não fornecido")
	}

	// Verificar formato "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.UnauthorizedResponse(c, "Formato de token inválido")
	}

	tokenString := parts[1]

	// Validar o token
	claims, err := utils.ValidateToken(tokenString, config.AppConfig.JWTSecret)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Token inválido ou expirado")
	}

	// Injetar userID no contexto para uso nos handlers
	c.Locals("userID", claims.UserID)

	return c.Next()
}

// GetUserID extrai o userID do contexto (helper function)
func GetUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}
