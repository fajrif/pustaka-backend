package models

import (
	"time"
	"github.com/google/uuid"
)

type SalesAssociate struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name             string     `gorm:"unique;not null" json:"name"`
	Address          string     `gorm:"not null" json:"address"`
	CityID           *uuid.UUID `gorm:"type:uuid" json:"city_id"`
	City             *City      `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Area             *string    `json:"area"`
	Phone1           string     `gorm:"not null" json:"phone1"`
	Phone2           *string    `json:"phone2"`
	Email            *string    `json:"email"`
	Website          *string    `json:"website"`
	JenisPembayaran  string     `gorm:"default:'T'" json:"jenis_pembayaran"`
	JoinDate         time.Time  `gorm:"not null" json:"join_date"`
	EndJoinDate      *time.Time `json:"end_join_date"`
	Discount         float64    `gorm:"not null" json:"discount"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (SalesAssociate) TableName() string {
	return "sales_associates"
}
