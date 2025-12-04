package models

import (
	"time"
	"github.com/google/uuid"
)

type BidangStudi struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BidangStudi) TableName() string {
	return "bidang_studi"
}
