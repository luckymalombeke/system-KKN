package services

import (
	"errors"
	"kkn-system/models/entity"
	"kkn-system/repositories"
	"strconv"
	"time"
)

type PembayaranService interface {
	CreateInvoice(pesertaID uint, amount int64) (entity.Pembayaran, error)
	GetStatus(id uint) (entity.Pembayaran, error)
}

type pembayaranService struct {
	repo        repositories.PembayaranRepository
	pesertaRepo repositories.PesertaRepository
}

func NewPembayaranService(repo repositories.PembayaranRepository, pesertaRepo repositories.PesertaRepository) PembayaranService {
	return &pembayaranService{repo, pesertaRepo}
}

func (s *pembayaranService) CreateInvoice(pesertaID uint, amount int64) (entity.Pembayaran, error) {
	// 1. Cek apakah peserta ada
	peserta, err := s.pesertaRepo.FindByID(pesertaID)
	if err != nil || peserta.ID == 0 {
		return entity.Pembayaran{}, errors.New("peserta tidak ditemukan")
	}

	// 2. Buat data pembayaran (simulasi Midtrans)
	externalID := "INV-" + strconv.FormatInt(time.Now().Unix(), 10)
	
	pembayaran := entity.Pembayaran{
		PesertaID:  pesertaID,
		Amount:     amount,
		Status:     "pending",
		ExternalID: externalID,
		PaymentURL: "https://checkout.midtrans.com/v2/vtweb/" + externalID, // Simulasi URL
	}

	return s.repo.Create(pembayaran)
}

func (s *pembayaranService) GetStatus(id uint) (entity.Pembayaran, error) {
	return s.repo.FindByID(id)
}
