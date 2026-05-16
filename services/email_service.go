package services

import "fmt"

type EmailService interface {
	SendEmail(to, subject, body string) error
}

type emailService struct {
	smtpHost string
}

func NewEmailService(smtpHost string) EmailService {
	return &emailService{smtpHost}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	// Logika kirim email (SMTP/SendGrid)
	fmt.Printf("Email terkirim ke %s: %s\n", to, subject)
	return nil
}
