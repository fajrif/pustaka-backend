package models

import (
	"time"
	"github.com/google/uuid"
)

type Expedition struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code        string     `gorm:"unique;not null" json:"code"`
	Name        string     `gorm:"unique;not null" json:"name"`
	Description *string    `json:"description"`
	Address     string     `gorm:"not null" json:"address"`
	CityID      *uuid.UUID `gorm:"type:uuid" json:"city_id"`
	City        *City      `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Area        *string    `json:"area"`
	Phone1      string     `gorm:"not null" json:"phone1"`
	Phone2      *string    `json:"phone2"`
	Email       *string    `json:"email"`
	Website     *string    `json:"website"`
	LogoUrl     *string    `json:"logo_url,omitempty"`
	FileUrl     *string    `json:"file_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Expedition) TableName() string {
	return "expeditions"
}
