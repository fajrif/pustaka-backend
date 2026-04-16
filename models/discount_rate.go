package models

import (
	"github.com/google/uuid"
	"time"
)

type DiscountRate struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Discount    float64    `gorm:"type:decimal(5,2);not null;default:0" json:"discount"`
	Periode     int        `gorm:"not null;default:1" json:"periode"`
	Year        string     `gorm:"not null" json:"year"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (DiscountRate) TableName() string {
	return "discount_rates"
}
