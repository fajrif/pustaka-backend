package models

import (
	"time"

	"github.com/google/uuid"
)

type Shipping struct {
	ID                 uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SalesTransactionID uuid.UUID   `gorm:"type:uuid;not null" json:"sales_transaction_id"`
	ExpeditionID       uuid.UUID   `gorm:"type:uuid;not null" json:"expedition_id"`
	Expedition         *Expedition `gorm:"foreignKey:ExpeditionID" json:"expedition,omitempty"`
	NoResi             *string     `json:"no_resi"`
	TotalAmount        float64     `gorm:"not null;default:0" json:"total_amount"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

func (Shipping) TableName() string {
	return "shippings"
}
