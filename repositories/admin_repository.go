package repositories

import (
	"kkn-system/models/entity"

	"gorm.io/gorm"
)

type AdminRepository interface {
	FindByEmail(email string) (entity.Admin, error)
	Create(admin entity.Admin) (entity.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db}
}

func (r *adminRepository) FindByEmail(email string) (entity.Admin, error) {
	var admin entity.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	return admin, err
}

func (r *adminRepository) Create(admin entity.Admin) (entity.Admin, error) {
	err := r.db.Create(&admin).Error
	return admin, err
}
