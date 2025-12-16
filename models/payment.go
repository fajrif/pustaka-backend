package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SalesTransactionID uuid.UUID `gorm:"type:uuid;not null" json:"sales_transaction_id"`
	NoPayment          string    `gorm:"unique;not null" json:"no_payment"`
	PaymentDate        time.Time `gorm:"not null" json:"payment_date"`
	Amount             float64   `gorm:"not null" json:"amount"`
	Note               *string   `json:"note"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}
