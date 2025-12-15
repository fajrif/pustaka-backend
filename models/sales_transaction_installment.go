package models

import (
	"time"
	"github.com/google/uuid"
)

type SalesTransactionInstallment struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TransactionID   uuid.UUID `gorm:"type:uuid;not null" json:"transaction_id"`
	NoInstallment   string    `gorm:"unique;not null" json:"no_installment"`
	InstallmentDate time.Time `gorm:"not null" json:"installment_date"`
	Amount          float64   `gorm:"not null" json:"amount"`
	Note            *string   `json:"note"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (SalesTransactionInstallment) TableName() string {
	return "sales_transaction_installments"
}
