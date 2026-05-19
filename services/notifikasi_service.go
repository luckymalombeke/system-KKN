package services

import (
	"errors"
	"kkn-system/models/entity"
	"kkn-system/repositories"
)

type NotifikasiService interface {
	Send(pesertaID uint, title, message string) (entity.Notifikasi, error)
	GetByPeserta(pesertaID uint) ([]entity.Notifikasi, error)
	Read(id uint) error
	ReadForPeserta(notifikasiID, pesertaID uint) error
}

type notifikasiService struct {
	repo repositories.NotifikasiRepository
}

func NewNotifikasiService(repo repositories.NotifikasiRepository) NotifikasiService {
	return &notifikasiService{repo}
}

func (s *notifikasiService) Send(pesertaID uint, title, message string) (entity.Notifikasi, error) {
	notifikasi := entity.Notifikasi{
		PesertaID: pesertaID,
		Title:     title,
		Message:   message,
	}
	return s.repo.Create(notifikasi)
}

func (s *notifikasiService) GetByPeserta(pesertaID uint) ([]entity.Notifikasi, error) {
	return s.repo.FindByPesertaID(pesertaID)
}

func (s *notifikasiService) Read(id uint) error {
	return s.repo.MarkAsRead(id)
}

func (s *notifikasiService) ReadForPeserta(notifikasiID, pesertaID uint) error {
	err := s.repo.MarkAsReadForPeserta(notifikasiID, pesertaID)
	if err != nil {
		return errors.New("notifikasi tidak ditemukan atau bukan milik Anda")
	}
	return nil
}
