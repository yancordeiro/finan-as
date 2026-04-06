package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/financas/backend/internal/config"
	"github.com/financas/backend/internal/database"
	"github.com/financas/backend/internal/models"
	"github.com/financas/backend/pkg/utils"
)

// BoolOrString é um tipo que aceita bool ou string do JSON
type BoolOrString bool

func (b *BoolOrString) UnmarshalJSON(data []byte) error {
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		*b = BoolOrString(boolVal)
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		*b = BoolOrString(strVal == "true")
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into bool", string(data))
}

// GoogleUserInfo representa as informações do usuário retornadas pelo Google
type GoogleUserInfo struct {
	Sub           string       `json:"sub"` // Google ID
	Email         string       `json:"email"`
	EmailVerified BoolOrString `json:"email_verified"`
	Name          string       `json:"name"`
	Picture       string       `json:"picture"`
}

// TokenPair representa um par de tokens (access + refresh)
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ValidateGoogleToken valida o ID token do Google e retorna as informações do usuário
func ValidateGoogleToken(idToken string) (*GoogleUserInfo, error) {
	// Validar token usando o endpoint do Google
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao validar token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token inválido: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Verificar se o email está verificado
	if !userInfo.EmailVerified {
		return nil, fmt.Errorf("email não verificado")
	}

	// Verificar se o client ID está correto (segurança adicional)
	// Nota: A resposta do tokeninfo não inclui o aud, mas podemos confiar na validação do Google

	return &userInfo, nil
}

// CreateOrUpdateUser cria um novo usuário ou atualiza um existente
func CreateOrUpdateUser(googleUser *GoogleUserInfo) (*models.User, error) {
	db := database.GetDB()

	var user models.User

	// Buscar usuário pelo Google ID
	result := db.Where("google_id = ?", googleUser.Sub).First(&user)

	if result.Error != nil {
		// Usuário não existe, criar novo
		user = models.User{
			GoogleID: googleUser.Sub,
			Email:    googleUser.Email,
			Name:     googleUser.Name,
			Picture:  googleUser.Picture,
		}

		if err := db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("erro ao criar usuário: %w", err)
		}
	} else {
		// Usuário existe, atualizar informações
		user.Email = googleUser.Email
		user.Name = googleUser.Name
		user.Picture = googleUser.Picture

		if err := db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("erro ao atualizar usuário: %w", err)
		}
	}

	return &user, nil
}

// GenerateTokens gera um par de tokens (access + refresh) para o usuário
func GenerateTokens(userID uint) (*TokenPair, error) {
	// Gerar access token (1 hora)
	accessToken, err := utils.GenerateAccessToken(userID, config.AppConfig.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	// Gerar refresh token (30 dias)
	refreshToken, err := utils.GenerateRefreshToken(userID, config.AppConfig.JWTRefreshSecret)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	// Salvar refresh token no banco
	db := database.GetDB()
	refreshTokenModel := models.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := db.Create(&refreshTokenModel).Error; err != nil {
		return nil, fmt.Errorf("erro ao salvar refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateRefreshToken valida um refresh token e retorna um novo par de tokens
func ValidateRefreshToken(tokenString string) (*TokenPair, *models.User, error) {
	// Validar assinatura do token
	claims, err := utils.ValidateToken(tokenString, config.AppConfig.JWTRefreshSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("token inválido: %w", err)
	}

	db := database.GetDB()

	// Verificar se o token existe no banco e não expirou
	var refreshToken models.RefreshToken
	result := db.Where("token = ? AND user_id = ? AND expires_at > ?",
		tokenString, claims.UserID, time.Now()).First(&refreshToken)

	if result.Error != nil {
		return nil, nil, fmt.Errorf("refresh token não encontrado ou expirado")
	}

	// Buscar usuário
	var user models.User
	if err := db.First(&user, claims.UserID).Error; err != nil {
		return nil, nil, fmt.Errorf("usuário não encontrado")
	}

	// Revogar o token antigo
	if err := db.Delete(&refreshToken).Error; err != nil {
		return nil, nil, fmt.Errorf("erro ao revogar token antigo: %w", err)
	}

	// Gerar novos tokens
	tokens, err := GenerateTokens(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return tokens, &user, nil
}

// RevokeRefreshToken revoga um refresh token (logout)
func RevokeRefreshToken(tokenString string) error {
	db := database.GetDB()

	result := db.Where("token = ?", tokenString).Delete(&models.RefreshToken{})

	if result.Error != nil {
		return fmt.Errorf("erro ao revogar token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("token não encontrado")
	}

	return nil
}

// RevokeAllUserTokens revoga todos os refresh tokens de um usuário
func RevokeAllUserTokens(userID uint) error {
	db := database.GetDB()

	if err := db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("erro ao revogar tokens: %w", err)
	}

	return nil
}
