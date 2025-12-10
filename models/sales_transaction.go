package models

import (
	"time"
	"github.com/google/uuid"
)

type SalesTransaction struct {
	ID                uuid.UUID                    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SalesAssociateID  uuid.UUID                    `gorm:"type:uuid;not null" json:"sales_associate_id"`
	SalesAssociate    *SalesAssociate              `gorm:"foreignKey:SalesAssociateID" json:"sales_associate,omitempty"`
	ExpeditionID      *uuid.UUID                   `gorm:"type:uuid" json:"expedition_id"`
	Expedition        *Expedition                  `gorm:"foreignKey:ExpeditionID" json:"expedition,omitempty"`
	NoInvoice         string                       `gorm:"unique;not null" json:"no_invoice"`
	PaymentType       string                       `gorm:"default:'T';not null" json:"payment_type"` // 'T' for Cash, 'K' for Credit
	TransactionDate   time.Time                    `gorm:"not null" json:"transaction_date"`
	DueDate           *time.Time                   `json:"due_date"`
	ExpeditionPrice   float64                      `gorm:"default:0" json:"expedition_price"`
	TotalAmount       float64                      `gorm:"not null;default:0" json:"total_amount"`
	Status            int                          `gorm:"not null;default:0" json:"status"` // 0 = booking, 1 = paid-off, 2 = installment
	Items             []SalesTransactionItem       `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Installments      []SalesTransactionInstallment `gorm:"foreignKey:TransactionID" json:"installments,omitempty"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
}

func (SalesTransaction) TableName() string {
	return "sales_transactions"
}
