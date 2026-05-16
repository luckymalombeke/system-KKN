package repositories

import (
	"kkn-system/models/entity"

	"gorm.io/gorm"
)

type PembayaranRepository interface {
	Create(pembayaran entity.Pembayaran) (entity.Pembayaran, error)
	FindByID(id uint) (entity.Pembayaran, error)
	FindByPesertaID(pesertaID uint) ([]entity.Pembayaran, error)
	UpdateStatus(externalID string, status string) error
}

type pembayaranRepository struct {
	db *gorm.DB
}

func NewPembayaranRepository(db *gorm.DB) PembayaranRepository {
	return &pembayaranRepository{db}
}

func (r *pembayaranRepository) Create(pembayaran entity.Pembayaran) (entity.Pembayaran, error) {
	err := r.db.Create(&pembayaran).Error
	return pembayaran, err
}

func (r *pembayaranRepository) FindByID(id uint) (entity.Pembayaran, error) {
	var pembayaran entity.Pembayaran
	err := r.db.Preload("Peserta").First(&pembayaran, id).Error
	return pembayaran, err
}

func (r *pembayaranRepository) FindByPesertaID(pesertaID uint) ([]entity.Pembayaran, error) {
	var pembayaran []entity.Pembayaran
	err := r.db.Where("peserta_id = ?", pesertaID).Find(&pembayaran).Error
	return pembayaran, err
}

func (r *pembayaranRepository) UpdateStatus(externalID string, status string) error {
	return r.db.Model(&entity.Pembayaran{}).Where("external_id = ?", externalID).Update("status", status).Error
}
