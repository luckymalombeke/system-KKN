package services

import (
	"errors"
	"fmt"
	"kkn-system/models/entity"
	"kkn-system/repositories"
	"kkn-system/utils"
	"strconv"
	"time"
)

type PembayaranService interface {
	CreateInvoice(pesertaID uint, amount int64) (entity.Pembayaran, error)
	GetStatus(id uint) (entity.Pembayaran, error)
	GetByPesertaID(pesertaID uint) (entity.Pembayaran, error)
	UpdateStatus(externalID string, status string) error
	HandleMidtransNotification(notification MidtransNotification) error
}

type pembayaranService struct {
	repo               repositories.PembayaranRepository
	pesertaRepo        repositories.PesertaRepository
	paymentService     PaymentService
	paymentExpiryHours int
}

func NewPembayaranService(
	repo repositories.PembayaranRepository,
	pesertaRepo repositories.PesertaRepository,
	paymentService PaymentService,
	paymentExpiryHours int,
) PembayaranService {
	if paymentExpiryHours <= 0 {
		paymentExpiryHours = 24
	}
	return &pembayaranService{repo, pesertaRepo, paymentService, paymentExpiryHours}
}

func (s *pembayaranService) CreateInvoice(pesertaID uint, amount int64) (entity.Pembayaran, error) {
	peserta, err := s.pesertaRepo.FindByID(pesertaID)
	if err != nil || peserta.ID == 0 {
		return entity.Pembayaran{}, errors.New("peserta tidak ditemukan")
	}

	if err := s.ensureCanCreateInvoice(pesertaID); err != nil {
		return entity.Pembayaran{}, err
	}

	externalID := "INV-" + strconv.FormatInt(time.Now().Unix(), 10)
	expiredAt := time.Now().Add(time.Duration(s.paymentExpiryHours) * time.Hour)

	redirectURL, err := s.paymentService.CreateSnapTransaction(externalID, amount, SnapCustomerDetail{
		FirstName: peserta.Nama,
		Email:     peserta.Email,
	}, expiredAt)
	if err != nil {
		return entity.Pembayaran{}, err
	}

	pembayaran := entity.Pembayaran{
		PesertaID:  pesertaID,
		Amount:     amount,
		Status:     "pending",
		ExternalID: externalID,
		PaymentURL: redirectURL,
		ExpiredAt:  expiredAt,
	}

	return s.repo.Create(pembayaran)
}

func (s *pembayaranService) ensureCanCreateInvoice(pesertaID uint) error {
	latest, err := s.repo.FindByPesertaID(pesertaID)
	if err != nil || latest.ID == 0 {
		return nil
	}

	latest = s.syncExpiredStatus(latest)

	switch latest.Status {
	case "success":
		return errors.New("pembayaran sudah lunas, tidak perlu membuat tagihan baru")
	case "pending":
		if !latest.ExpiredAt.IsZero() && time.Now().After(latest.ExpiredAt) {
			_ = s.repo.UpdateStatus(latest.ExternalID, "failed")
			return nil
		}
		return fmt.Errorf("masih ada tagihan aktif. Bayar sebelum %s atau tunggu kedaluwarsa", formatWIB(latest.ExpiredAt))
	default:
		return nil
	}
}

func (s *pembayaranService) syncExpiredStatus(p entity.Pembayaran) entity.Pembayaran {
	if p.Status == "pending" && !p.ExpiredAt.IsZero() && time.Now().After(p.ExpiredAt) {
		_ = s.repo.UpdateStatus(p.ExternalID, "failed")
		p.Status = "failed"
	}
	return p
}

func formatWIB(t time.Time) string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60)
	}
	return t.In(loc).Format("02 Jan 2006, 15:04 WIB")
}

func (s *pembayaranService) GetStatus(id uint) (entity.Pembayaran, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return p, err
	}
	return s.syncExpiredStatus(p), nil
}

func (s *pembayaranService) GetByPesertaID(pesertaID uint) (entity.Pembayaran, error) {
	p, err := s.repo.FindByPesertaID(pesertaID)
	if err != nil {
		return p, err
	}
	return s.syncExpiredStatus(p), nil
}

func (s *pembayaranService) UpdateStatus(externalID string, status string) error {
	return s.repo.UpdateStatus(externalID, status)
}

func (s *pembayaranService) HandleMidtransNotification(notification MidtransNotification) error {
	if notification.OrderID == "" || notification.TransactionStatus == "" {
		return errors.New("order_id dan transaction_status wajib ada")
	}

	if err := s.paymentService.VerifyNotification(notification); err != nil {
		return err
	}

	pembayaran, err := s.repo.FindByExternalID(notification.OrderID)
	if err != nil || pembayaran.ID == 0 {
		return errors.New("order_id tidak ditemukan")
	}

	if !utils.GrossAmountMatches(notification.GrossAmount, pembayaran.Amount) {
		return errors.New("gross_amount tidak sesuai dengan tagihan")
	}

	finalStatus := MapMidtransTransactionStatus(notification)
	return s.repo.UpdateStatus(notification.OrderID, finalStatus)
}
