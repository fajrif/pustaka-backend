package models

import (
	"time"
	"github.com/google/uuid"
)

type Kelas struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code        string    `gorm:"unique;not null" json:"code"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Kelas) TableName() string {
	return "kelas"
}
