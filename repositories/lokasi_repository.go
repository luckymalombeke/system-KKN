package repositories

import (
	"kkn-system/models/entity"

	"gorm.io/gorm"
)

type LokasiRepository interface {
	Create(lokasi entity.Lokasi) (entity.Lokasi, error)
	FindAll() ([]entity.Lokasi, error)
	FindByID(id uint) (entity.Lokasi, error)
}

type lokasiRepository struct {
	db *gorm.DB
}

func NewLokasiRepository(db *gorm.DB) LokasiRepository {
	return &lokasiRepository{db}
}

func (r *lokasiRepository) Create(lokasi entity.Lokasi) (entity.Lokasi, error) {
	err := r.db.Create(&lokasi).Error
	return lokasi, err
}

func (r *lokasiRepository) FindAll() ([]entity.Lokasi, error) {
	var lokasi []entity.Lokasi
	err := r.db.Find(&lokasi).Error
	return lokasi, err
}

func (r *lokasiRepository) FindByID(id uint) (entity.Lokasi, error) {
	var lokasi entity.Lokasi
	err := r.db.First(&lokasi, id).Error
	return lokasi, err
}
