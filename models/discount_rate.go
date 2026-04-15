package models

import (
	"github.com/google/uuid"
	"time"
)

type DiscountRate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Discount    float64   `gorm:"type:decimal(5,2);not null;default:0" json:"discount"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (DiscountRate) TableName() string {
	return "discount_rates"
}
