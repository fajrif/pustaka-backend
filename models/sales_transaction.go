package models

import (
	"time"
	"github.com/google/uuid"
)

type SalesTransaction struct {
	ID               uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BillerID         *uuid.UUID             `gorm:"type:uuid" json:"biller_id"`
	Biller           *Biller                `gorm:"foreignKey:BillerID" json:"biller,omitempty"`
	SalesAssociateID uuid.UUID              `gorm:"type:uuid;not null" json:"sales_associate_id"`
	SalesAssociate   *SalesAssociate        `gorm:"foreignKey:SalesAssociateID" json:"sales_associate,omitempty"`
	NoInvoice        string                 `gorm:"unique;not null" json:"no_invoice"`
	PaymentType      string                 `gorm:"default:'T';not null" json:"payment_type"` // 'T' for Cash, 'K' for Credit
	TransactionDate  time.Time              `gorm:"not null" json:"transaction_date"`
	DueDate          *time.Time             `json:"due_date"`
	TotalAmount      float64                `gorm:"not null;default:0" json:"total_amount"`
	Status           int                    `gorm:"not null;default:0" json:"status"` // 0 = booking, 1 = paid-off, 2 = installment
	Items            []SalesTransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Payments         []Payment              `gorm:"foreignKey:SalesTransactionID" json:"payments,omitempty"`
	Shippings        []Shipping             `gorm:"foreignKey:SalesTransactionID" json:"shippings,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

func (SalesTransaction) TableName() string {
	return "sales_transactions"
}
