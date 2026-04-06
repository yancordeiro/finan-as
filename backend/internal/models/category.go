package models

import (
	"time"

	"gorm.io/gorm"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
	CategoryTypeBoth    CategoryType = "both"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    *uint          `gorm:"index" json:"userId"` // NULL para categorias globais
	Name      string         `gorm:"not null" json:"name"`
	Icon      string         `gorm:"not null" json:"icon"` // Emoji ou nome do ícone
	Color     string         `gorm:"not null" json:"color"` // Hex color
	Type      CategoryType   `gorm:"type:varchar(20);not null;default:'expense'" json:"type"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         *User         `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID" json:"transactions,omitempty"`
}

// IsGlobal verifica se a categoria é global (padrão do sistema)
func (c *Category) IsGlobal() bool {
	return c.UserID == nil
}
