package entity

import (
	"time"

	"gorm.io/gorm"
)

type Peserta struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"type:varchar(255);not null" json:"nama"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	NIM       string         `gorm:"type:varchar(50);unique;not null" json:"nim"`
	Prodi     string         `gorm:"type:varchar(100)" json:"prodi"`
	Role      string         `gorm:"type:varchar(20);default:'mahasiswa'" json:"role"`
	Status    string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
	LokasiID     *uint          `json:"lokasi_id"`
	Lokasi       *Lokasi        `gorm:"foreignKey:LokasiID" json:"lokasi,omitempty"`
	OTP          string         `gorm:"type:varchar(10)" json:"-"` // Tidak dikembalikan di JSON
	OTPExpiredAt time.Time      `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
