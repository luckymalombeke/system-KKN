package services

import (
	"errors"
	"fmt"
	"kkn-system/config"
	"log"
	"net/smtp"
	"strings"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendOTPEmail(to, nama, otp string, validMinutes int) error
}

type emailService struct {
	cfg config.SMTPConfig
}

func NewEmailService(cfg config.SMTPConfig) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendOTPEmail(to, nama, otp string, validMinutes int) error {
	subject := "Kode OTP Login KKN"
	body := fmt.Sprintf(
		"Halo %s,\n\nKode OTP login KKN Anda: %s\nKode berlaku %d menit.\n\nJangan bagikan kode ini kepada siapa pun.\n",
		nama, otp, validMinutes,
	)
	return s.SendEmail(to, subject, body)
}

func (s *emailService) SendEmail(to, subject, body string) error {
	to = strings.TrimSpace(to)
	if to == "" {
		return errors.New("alamat email penerima kosong")
	}

	if s.cfg.DevLog {
		log.Printf("[OTP_DEV_LOG] email ke %s | subjek: %s | isi:\n%s", to, subject, body)
		return nil
	}

	if s.cfg.Host == "" {
		return errors.New("SMTP belum dikonfigurasi (isi SMTP_HOST atau aktifkan OTP_DEV_LOG=true untuk pengujian lokal)")
	}

	from := s.cfg.From
	if from == "" {
		from = s.cfg.User
	}
	if from == "" {
		return errors.New("SMTP_FROM atau SMTP_USER wajib diisi")
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	msg := buildPlainTextMessage(from, to, subject, body)

	var auth smtp.Auth
	if s.cfg.User != "" {
		auth = smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	}

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	return nil
}

func buildPlainTextMessage(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	)
}
