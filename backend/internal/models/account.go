package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeInvestment AccountType = "investment"
)

type Account struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"userId"`
	Name      string         `gorm:"not null" json:"name"`
	Type      AccountType    `gorm:"type:varchar(20);not null;default:'checking'" json:"type"`
	Balance   int64          `gorm:"default:0" json:"balance"` // Valor em centavos
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         User          `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:AccountID" json:"transactions,omitempty"`
}

// BalanceInReais retorna o saldo em reais (float64)
func (a *Account) BalanceInReais() float64 {
	return float64(a.Balance) / 100.0
}
