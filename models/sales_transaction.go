package models

import (
	"time"
	"github.com/google/uuid"
)

type SalesTransaction struct {
	ID               uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BillerID         *uuid.UUID             `gorm:"type:uuid" json:"biller_id"`
	Biller           *Biller                `gorm:"foreignKey:BillerID" json:"biller,omitempty"`
	SalesAssociateID uuid.UUID              `gorm:"type:uuid;not null" json:"sales_associate_id"`
	SalesAssociate   *SalesAssociate        `gorm:"foreignKey:SalesAssociateID" json:"sales_associate,omitempty"`
	NoInvoice        string                 `gorm:"unique;not null" json:"no_invoice"`
	PaymentType      string                 `gorm:"default:'T';not null" json:"payment_type"` // 'T' for Cash, 'K' for Credit
	TransactionDate  time.Time              `gorm:"not null" json:"transaction_date"`
	DueDate          *time.Time             `json:"due_date"`
	SecondaryDueDate *time.Time             `json:"secondary_due_date"`
	TotalAmount      float64                `gorm:"not null;default:0" json:"total_amount"`
	Status           int                    `gorm:"not null;default:0" json:"status"` // 0 = booking, 1 = paid-off, 2 = installment
	Periode          int                    `gorm:"not null;default:1" json:"periode"`
	Year             string                 `gorm:"not null" json:"year"`
	CurriculumID     *uuid.UUID             `gorm:"type:uuid" json:"curriculum_id"`
	Curriculum       *Curriculum            `gorm:"foreignKey:CurriculumID" json:"curriculum,omitempty"`
	MerkBukuID       *uuid.UUID             `gorm:"type:uuid" json:"merk_buku_id"`
	MerkBuku         *MerkBuku              `gorm:"foreignKey:MerkBukuID" json:"merk_buku,omitempty"`
	JenjangStudiID   *uuid.UUID             `gorm:"type:uuid" json:"jenjang_studi_id"`
	JenjangStudi     *JenjangStudi          `gorm:"foreignKey:JenjangStudiID" json:"jenjang_studi,omitempty"`
	JenisBukuID      *uuid.UUID             `gorm:"type:uuid" json:"jenis_buku_id"`
	JenisBuku        *JenisBuku             `gorm:"foreignKey:JenisBukuID" json:"jenis_buku,omitempty"`
	Items            []SalesTransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Payments         []Payment              `gorm:"foreignKey:SalesTransactionID" json:"payments,omitempty"`
	Shippings        []Shipping             `gorm:"foreignKey:SalesTransactionID" json:"shippings,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

func (SalesTransaction) TableName() string {
	return "sales_transactions"
}
