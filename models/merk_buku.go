package models

import (
	"time"
	"github.com/google/uuid"
)

// MerkBuku represents a book brand/publisher
type MerkBuku struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KodeMerk       *string   `gorm:"column:kode_merk;not null" json:"kode_merk"`
	NamaMerk       *string   `gorm:"column:nama_merk;not null" json:"nama_merk"`
	BantuanPromosi *int      `json:"bantuan_promosi"`
	UserID         uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User           User 		 `gorm:"foreignKey:UserID" json:"user_details"`
	Tstamp         time.Time `json:"tstamp" db:"tstamp"`
}

func (MerkBuku) TableName() string {
    return "merk_buku"
}
