package models

import (
	"time"
	"github.com/google/uuid"
)

type JenisBuku struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (JenisBuku) TableName() string {
	return "jenis_buku"
}
