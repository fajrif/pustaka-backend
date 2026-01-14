package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseTransactionItem struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PurchaseTransactionID uuid.UUID `gorm:"type:uuid;not null" json:"purchase_transaction_id"`
	BookID                uuid.UUID `gorm:"type:uuid;not null" json:"book_id"`
	Book                  *Book     `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity              int       `gorm:"not null" json:"quantity"`
	Price                 float64   `gorm:"not null" json:"price"`
	Subtotal              float64   `gorm:"not null" json:"subtotal"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func (PurchaseTransactionItem) TableName() string {
	return "purchase_transaction_items"
}
