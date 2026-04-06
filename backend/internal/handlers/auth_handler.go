package handlers

import (
	"log"
	"time"

	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/middleware"
	"github.com/financas/backend/internal/services"
	"github.com/financas/backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// GoogleAuthRequest representa o corpo da requisição de login com Google
type GoogleAuthRequest struct {
	IDToken string `json:"idToken"`
}

// GoogleAuthResponse representa a resposta do login
type GoogleAuthResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         interface{} `json:"user"`
}

// GoogleAuth handler para autenticação com Google OAuth
func GoogleAuth(c *fiber.Ctx) error {
	var req GoogleAuthRequest

	// Parse do body
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Dados inválidos")
	}

	// Validar que o idToken foi enviado
	if req.IDToken == "" {
		return utils.ValidationErrorResponse(c, "ID Token é obrigatório")
	}

	// Validar o token com o Google
	googleUser, err := services.ValidateGoogleToken(req.IDToken)
	if err != nil {
		log.Printf("Erro ao validar token do Google: %v", err)
		return utils.UnauthorizedResponse(c, "Token do Google inválido")
	}

	// Criar ou atualizar usuário no banco
	user, err := services.CreateOrUpdateUser(googleUser)
	if err != nil {
		log.Printf("Erro ao criar/atualizar usuário: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao processar usuário")
	}

	// Gerar tokens JWT
	tokens, err := services.GenerateTokens(user.ID)
	if err != nil {
		log.Printf("Erro ao gerar tokens: %v", err)
		return utils.InternalErrorResponse(c, "Erro ao gerar tokens")
	}

	// Definir refresh token em cookie httpOnly
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Mudar para true em produção (HTTPS)
		SameSite: "Lax",
	})

	// Retornar resposta
	return utils.SuccessResponse(c, fiber.StatusOK, "Login realizado com sucesso", GoogleAuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
		},
	})
}

// RefreshToken handler para renovar o access token
func RefreshToken(c *fiber.Ctx) error {
	// Ler refresh token do cookie
	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" {
		return utils.UnauthorizedResponse(c, "Refresh token não fornecido")
	}

	// Validar refresh token e gerar novos tokens
	tokens, user, err := services.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Erro ao validar refresh token: %v", err)
		return utils.UnauthorizedResponse(c, "Refresh token inválido ou expirado")
	}

	// Atualizar cookie com novo refresh token
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Mudar para true em produção
		SameSite: "Lax",
	})

	// Retornar novos tokens
	return utils.SuccessResponse(c, fiber.StatusOK, "Token renovado com sucesso", GoogleAuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User: map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"picture": user.Picture,
		},
	})
}

// Logout handler para fazer logout (revogar refresh token)
func Logout(c *fiber.Ctx) error {
	// Ler refresh token do cookie
	refreshToken := c.Cookies("refreshToken")

	if refreshToken != "" {
		// Revogar o refresh token
		if err := services.RevokeRefreshToken(refreshToken); err != nil {
			log.Printf("Erro ao revogar refresh token: %v", err)
			// Não retornar erro, continuar com o logout
		}
	}

	// Limpar cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	return utils.SuccessResponse(c, fiber.StatusOK, "Logout realizado com sucesso", nil)
}

// GetCurrentUser retorna as informações do usuário autenticado
func GetCurrentUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	if userID == 0 {
		return utils.UnauthorizedResponse(c, "Usuário não autenticado")
	}

	// Buscar usuário no banco
	var user struct {
		ID      uint   `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	db := database.GetDB()
	if err := db.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		log.Printf("Erro ao buscar usuário: %v", err)
		return utils.NotFoundResponse(c, "Usuário não encontrado")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Usuário encontrado", user)
}
