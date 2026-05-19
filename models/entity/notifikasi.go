package entity

import (
	"time"

	"gorm.io/gorm"
)

type Notifikasi struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PesertaID uint           `gorm:"not null" json:"peserta_id"`
	Peserta   Peserta        `gorm:"foreignKey:PesertaID" json:"peserta,omitempty"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
