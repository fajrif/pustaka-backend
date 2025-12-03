package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    Email        string    `gorm:"unique;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"`
    FullName     string    `gorm:"not null" json:"full_name"`
    Role         string    `gorm:"default:'user'" json:"role"`
    CreatedAt    time.Time `json:"created_date"`
    UpdatedAt    time.Time `json:"updated_date"`
}

func (User) TableName() string {
    return "users"
}
