package entity

import (
	"time"

	"gorm.io/gorm"
)

type Pembayaran struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	PesertaID     uint           `gorm:"not null" json:"peserta_id"`
	Peserta       Peserta        `gorm:"foreignKey:PesertaID" json:"peserta,omitempty"`
	Amount        int64          `gorm:"not null" json:"amount"`
	Status        string         `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, success, failed
	ExternalID    string         `gorm:"type:varchar(255)" json:"external_id"`             // ID dari Midtrans
	PaymentURL    string         `gorm:"type:text" json:"payment_url"`                     // URL Snap Midtrans
	ExpiredAt     time.Time      `json:"expired_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
