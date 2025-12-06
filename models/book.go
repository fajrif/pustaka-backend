package models

import (
	"time"
	"github.com/google/uuid"
)

type Book struct {
	ID             uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name           string        `gorm:"not null" json:"name"`
	Description    *string       `json:"description"`
	Year           string        `gorm:"not null" json:"year"`
	Author         *string       `json:"author"`
	ISBN           *string       `json:"isbn"`
	Stock          int           `gorm:"default:0" json:"stock"`
	JenisBukuID    *uuid.UUID    `gorm:"type:uuid" json:"jenis_buku_id"`
	JenisBuku      *JenisBuku    `gorm:"foreignKey:JenisBukuID" json:"jenis_buku,omitempty"`
	JenjangStudiID *uuid.UUID    `gorm:"type:uuid" json:"jenjang_studi_id"`
	JenjangStudi   *JenjangStudi `gorm:"foreignKey:JenjangStudiID" json:"jenjang_studi,omitempty"`
	BidangStudiID  *uuid.UUID    `gorm:"type:uuid" json:"bidang_studi_id"`
	BidangStudi    *BidangStudi  `gorm:"foreignKey:BidangStudiID" json:"bidang_studi,omitempty"`
	KelasID        *uuid.UUID    `gorm:"type:uuid" json:"kelas_id"`
	Kelas          *Kelas        `gorm:"foreignKey:KelasID" json:"kelas,omitempty"`
	PublisherID    *uuid.UUID    `gorm:"type:uuid" json:"publisher_id"`
	Publisher      *Publisher    `gorm:"foreignKey:PublisherID" json:"publisher,omitempty"`
	Price          float64       `gorm:"not null" json:"price"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func (Book) TableName() string {
	return "books"
}
