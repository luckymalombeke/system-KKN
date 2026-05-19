package services

import (
	"errors"
	"kkn-system/models/entity"
	"kkn-system/repositories"
	"kkn-system/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RequestOTP(nim string) error
	VerifyOTP(nim string, otp string) (string, entity.Peserta, error)
	AdminLogin(email string, password string) (string, entity.Admin, error)
}

type authService struct {
	pesertaRepo      repositories.PesertaRepository
	adminRepo        repositories.AdminRepository
	emailService     EmailService
	otpExpiryMinutes int
}

func NewAuthService(
	pesertaRepo repositories.PesertaRepository,
	adminRepo repositories.AdminRepository,
	emailService EmailService,
	otpExpiryMinutes int,
) AuthService {
	if otpExpiryMinutes <= 0 {
		otpExpiryMinutes = 5
	}
	return &authService{
		pesertaRepo:      pesertaRepo,
		adminRepo:        adminRepo,
		emailService:     emailService,
		otpExpiryMinutes: otpExpiryMinutes,
	}
}

func (s *authService) RequestOTP(nim string) error {
	nim = strings.TrimSpace(nim)
	peserta, err := s.pesertaRepo.FindByNIM(nim)
	if err != nil || peserta.ID == 0 {
		return errors.New("NIM tidak terdaftar, silakan daftar KKN terlebih dahulu")
	}

	otp, err := utils.GenerateOTP6()
	if err != nil {
		return errors.New("gagal membuat OTP")
	}

	expiredAt := time.Now().Add(time.Duration(s.otpExpiryMinutes) * time.Minute)

	if err := s.emailService.SendOTPEmail(peserta.Email, peserta.Nama, otp, s.otpExpiryMinutes); err != nil {
		return err
	}

	if err := s.pesertaRepo.UpdateOTP(peserta.ID, otp, expiredAt); err != nil {
		return errors.New("gagal menyimpan OTP")
	}

	return nil
}

func (s *authService) VerifyOTP(nim string, inputOTP string) (string, entity.Peserta, error) {
	nim = strings.TrimSpace(nim)
	inputOTP = strings.TrimSpace(inputOTP)

	peserta, err := s.pesertaRepo.FindByNIM(nim)
	if err != nil || peserta.ID == 0 {
		return "", entity.Peserta{}, errors.New("NIM tidak ditemukan")
	}

	if peserta.OTP != inputOTP {
		return "", entity.Peserta{}, errors.New("kode OTP salah")
	}

	if time.Now().After(peserta.OTPExpiredAt) {
		return "", entity.Peserta{}, errors.New("kode OTP sudah kedaluwarsa")
	}

	_ = s.pesertaRepo.UpdateOTP(peserta.ID, "", time.Now())

	token, err := utils.GenerateToken(peserta.ID, peserta.Role)
	if err != nil {
		return "", entity.Peserta{}, errors.New("gagal membuat token auth")
	}

	return token, peserta, nil
}

func (s *authService) AdminLogin(email string, password string) (string, entity.Admin, error) {
	admin, err := s.adminRepo.FindByEmail(email)
	if err != nil || admin.ID == 0 {
		return "", entity.Admin{}, errors.New("email admin tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", entity.Admin{}, errors.New("password salah")
	}

	token, err := utils.GenerateToken(admin.ID, admin.Role)
	if err != nil {
		return "", entity.Admin{}, errors.New("gagal membuat token admin")
	}

	return token, admin, nil
}
