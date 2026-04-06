package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	UserID      uint            `gorm:"index;not null" json:"userId"`
	CategoryID  *uint           `gorm:"index" json:"categoryId"` // Nullable - permite transações sem categoria
	Type        TransactionType `gorm:"type:varchar(20);not null" json:"type"`
	Description string          `gorm:"not null" json:"description"`
	Amount      int64           `gorm:"not null" json:"amount"` // Valor em centavos, sempre positivo
	Date        time.Time       `gorm:"index;not null" json:"date"`
	Notes       string          `json:"notes"`
	OFXID       string          `gorm:"index" json:"ofxId"` // FITID do arquivo OFX (para deduplicação)
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// AmountInReais retorna o valor em reais (float64)
func (t *Transaction) AmountInReais() float64 {
	return float64(t.Amount) / 100.0
}

// AmountSigned retorna o valor com sinal (negativo para expense, positivo para income)
func (t *Transaction) AmountSigned() int64 {
	if t.Type == TransactionTypeExpense {
		return -t.Amount
	}
	return t.Amount
}

// AmountSignedInReais retorna o valor com sinal em reais
func (t *Transaction) AmountSignedInReais() float64 {
	return float64(t.AmountSigned()) / 100.0
}
