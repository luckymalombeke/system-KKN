package repositories

import (
	"errors"
	"kkn-system/models/entity"

	"gorm.io/gorm"
)

type NotifikasiRepository interface {
	Create(notifikasi entity.Notifikasi) (entity.Notifikasi, error)
	FindByPesertaID(pesertaID uint) ([]entity.Notifikasi, error)
	MarkAsRead(id uint) error
	MarkAsReadForPeserta(notifikasiID, pesertaID uint) error
}

type notifikasiRepository struct {
	db *gorm.DB
}

func NewNotifikasiRepository(db *gorm.DB) NotifikasiRepository {
	return &notifikasiRepository{db}
}

func (r *notifikasiRepository) Create(notifikasi entity.Notifikasi) (entity.Notifikasi, error) {
	err := r.db.Create(&notifikasi).Error
	return notifikasi, err
}

func (r *notifikasiRepository) FindByPesertaID(pesertaID uint) ([]entity.Notifikasi, error) {
	var notifikasi []entity.Notifikasi
	err := r.db.Where("peserta_id = ?", pesertaID).Order("created_at desc").Find(&notifikasi).Error
	return notifikasi, err
}

func (r *notifikasiRepository) MarkAsRead(id uint) error {
	return r.db.Model(&entity.Notifikasi{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notifikasiRepository) MarkAsReadForPeserta(notifikasiID, pesertaID uint) error {
	result := r.db.Model(&entity.Notifikasi{}).
		Where("id = ? AND peserta_id = ?", notifikasiID, pesertaID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("notifikasi tidak ditemukan")
	}
	return nil
}
