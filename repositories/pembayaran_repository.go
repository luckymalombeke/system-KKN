package repositories

import (
	"errors"
	"kkn-system/models/entity"

	"gorm.io/gorm"
)

type PembayaranRepository interface {
	Create(pembayaran entity.Pembayaran) (entity.Pembayaran, error)
	FindByID(id uint) (entity.Pembayaran, error)
	FindByPesertaID(pesertaID uint) (entity.Pembayaran, error)
	FindByExternalID(externalID string) (entity.Pembayaran, error)
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

func (r *pembayaranRepository) FindByPesertaID(pesertaID uint) (entity.Pembayaran, error) {
	var pembayaran entity.Pembayaran
	err := r.db.Where("peserta_id = ?", pesertaID).Order("created_at desc").First(&pembayaran).Error
	return pembayaran, err
}

func (r *pembayaranRepository) FindByExternalID(externalID string) (entity.Pembayaran, error) {
	var pembayaran entity.Pembayaran
	err := r.db.Where("external_id = ?", externalID).First(&pembayaran).Error
	return pembayaran, err
}

func (r *pembayaranRepository) UpdateStatus(externalID string, status string) error {
	result := r.db.Model(&entity.Pembayaran{}).Where("external_id = ?", externalID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("data pembayaran dengan order_id tersebut tidak ditemukan")
	}
	return nil
}
