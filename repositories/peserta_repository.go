package repositories

import (
	"kkn-system/models/entity"

	"gorm.io/gorm"
)

type PesertaRepository interface {
	Create(peserta entity.Peserta) (entity.Peserta, error)
	FindAll() ([]entity.Peserta, error)
	FindByID(id uint) (entity.Peserta, error)
	FindByNIM(nim string) (entity.Peserta, error)
	Update(peserta entity.Peserta) (entity.Peserta, error)
}

type pesertaRepository struct {
	db *gorm.DB
}

func NewPesertaRepository(db *gorm.DB) PesertaRepository {
	return &pesertaRepository{db}
}

func (r *pesertaRepository) Create(peserta entity.Peserta) (entity.Peserta, error) {
	err := r.db.Create(&peserta).Error
	return peserta, err
}

func (r *pesertaRepository) FindAll() ([]entity.Peserta, error) {
	var peserta []entity.Peserta
	err := r.db.Find(&peserta).Error
	return peserta, err
}

func (r *pesertaRepository) FindByID(id uint) (entity.Peserta, error) {
	var peserta entity.Peserta
	err := r.db.First(&peserta, id).Error
	return peserta, err
}

func (r *pesertaRepository) FindByNIM(nim string) (entity.Peserta, error) {
	var peserta entity.Peserta
	err := r.db.Where("nim = ?", nim).First(&peserta).Error
	return peserta, err
}

func (r *pesertaRepository) Update(peserta entity.Peserta) (entity.Peserta, error) {
	err := r.db.Save(&peserta).Error
	return peserta, err
}
