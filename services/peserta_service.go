package services

import (
	"errors"
	"kkn-system/models/entity"
	"kkn-system/repositories"
)

type PesertaService interface {
	DaftarKKN(peserta entity.Peserta) (entity.Peserta, error)
	GetAllPeserta() ([]entity.Peserta, error)
	GetPesertaByID(id uint) (entity.Peserta, error)
	UpdateStatus(id uint, status string) error
}

type pesertaService struct {
	repo repositories.PesertaRepository
}

func NewPesertaService(repo repositories.PesertaRepository) PesertaService {
	return &pesertaService{repo}
}

func (s *pesertaService) DaftarKKN(peserta entity.Peserta) (entity.Peserta, error) {
	// Business Logic: Check if NIM already exists
	existing, _ := s.repo.FindByNIM(peserta.NIM)
	if existing.ID != 0 {
		return entity.Peserta{}, errors.New("NIM sudah terdaftar")
	}

	peserta.Status = "pending"
	return s.repo.Create(peserta)
}

func (s *pesertaService) GetAllPeserta() ([]entity.Peserta, error) {
	return s.repo.FindAll()
}

func (s *pesertaService) GetPesertaByID(id uint) (entity.Peserta, error) {
	return s.repo.FindByID(id)
}

func (s *pesertaService) UpdateStatus(id uint, status string) error {
	// Validasi: hanya status tertentu yang diperbolehkan
	validStatus := map[string]bool{"pending": true, "approved": true, "rejected": true}
	if !validStatus[status] {
		return errors.New("status tidak valid (hanya: pending, approved, rejected)")
	}

	return s.repo.UpdateStatus(id, status)
}
