package services

import (
	"kkn-system/models/entity"
	"kkn-system/repositories"
)

type LokasiService interface {
	CreateLokasi(lokasi entity.Lokasi) (entity.Lokasi, error)
	GetAllLokasi() ([]entity.Lokasi, error)
}

type lokasiService struct {
	repo repositories.LokasiRepository
}

func NewLokasiService(repo repositories.LokasiRepository) LokasiService {
	return &lokasiService{repo}
}

func (s *lokasiService) CreateLokasi(lokasi entity.Lokasi) (entity.Lokasi, error) {
	return s.repo.Create(lokasi)
}

func (s *lokasiService) GetAllLokasi() ([]entity.Lokasi, error) {
	return s.repo.FindAll()
}
