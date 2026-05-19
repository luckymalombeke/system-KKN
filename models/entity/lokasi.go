package entity

import (
	"time"

	"gorm.io/gorm"
)

type Lokasi struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NamaDesa  string         `gorm:"type:varchar(255);not null" json:"nama_desa"`
	Kecamatan string         `gorm:"type:varchar(255);not null" json:"kecamatan"`
	Kuota     int            `gorm:"default:0" json:"kuota"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
