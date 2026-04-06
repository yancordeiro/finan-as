package database

import (
	"fmt"
	"log"

	"github.com/financas/backend/internal/config"
	"github.com/financas/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect estabelece conexão com o banco de dados PostgreSQL
func Connect() error {
	var err error

	// Configuração do GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Conectar ao PostgreSQL
	DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), gormConfig)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	log.Println("✓ Conectado ao banco de dados PostgreSQL")

	// Executar migrations
	if err := RunMigrations(); err != nil {
		return fmt.Errorf("erro ao executar migrations: %w", err)
	}

	// Executar seeds
	if err := SeedDefaultCategories(); err != nil {
		return fmt.Errorf("erro ao executar seed de categorias: %w", err)
	}

	return nil
}

// RunMigrations executa as migrations automáticas do GORM
func RunMigrations() error {
	log.Println("Executando migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Category{},
		&models.Transaction{},
	)

	if err != nil {
		return err
	}

	log.Println("✓ Migrations executadas com sucesso")
	return nil
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	return DB
}
