package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseTransaction struct {
	ID              uuid.UUID                 `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SupplierID      uuid.UUID                 `gorm:"type:uuid;not null" json:"supplier_id"`
	Supplier        *Publisher                `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	NoInvoice       string                    `gorm:"unique;not null" json:"no_invoice"`
	PurchaseDate    time.Time                 `gorm:"not null" json:"purchase_date"`
	TotalAmount     float64                   `gorm:"not null;default:0" json:"total_amount"`
	Status          int                       `gorm:"not null;default:0" json:"status"` // 0 = pending, 1 = completed, 2 = cancelled
	ReceiptImageUrl *string                   `json:"receipt_image_url"`
	Note            *string                   `json:"note"`
	Items           []PurchaseTransactionItem `gorm:"foreignKey:PurchaseTransactionID" json:"items,omitempty"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

func (PurchaseTransaction) TableName() string {
	return "purchase_transactions"
}

// Status constants for PurchaseTransaction
const (
	PurchaseStatusPending   = 0 // Draft, stock not affected
	PurchaseStatusCompleted = 1 // Stock increased
	PurchaseStatusCancelled = 2 // Cancelled, stock not affected
)
