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
