package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTRefreshSecret   string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	FrontendURL        string
	Port               string
}

var AppConfig *Config

// LoadConfig carrega as variáveis de ambiente e inicializa a configuração
func LoadConfig() {
	// Tenta carregar .env file (ignorar erro se não existir, pois pode usar env vars do sistema)
	_ = godotenv.Load()

	AppConfig = &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://financas:financas_dev_password@localhost:5432/financas_db?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURI:  getEnv("GOOGLE_REDIRECT_URI", "http://localhost:3000/auth/callback"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		Port:               getEnv("PORT", "8080"),
	}

	// Validar configurações obrigatórias
	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET é obrigatório")
	}
	if AppConfig.JWTRefreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET é obrigatório")
	}
	if AppConfig.GoogleClientID == "" {
		log.Println("AVISO: GOOGLE_CLIENT_ID não configurado - autenticação Google não funcionará")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
