package models

import (
	"time"
	"github.com/google/uuid"
)

// MerkBuku represents a book brand
type MerkBuku struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code        string    `gorm:"unique;not null" json:"code"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description *string   `json:"description"`
	BantuanPromosi *int      `json:"bantuan_promosi"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MerkBuku) TableName() string {
    return "merk_buku"
}
