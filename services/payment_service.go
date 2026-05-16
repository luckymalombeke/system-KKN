package services

type PaymentService interface {
	CreateTransaction(amount int64, orderID string) (string, error)
}

type paymentService struct {
	serverKey string
}

func NewPaymentService(serverKey string) PaymentService {
	return &paymentService{serverKey}
}

func (s *paymentService) CreateTransaction(amount int64, orderID string) (string, error) {
	// Integrasi Midtrans di sini
	return "https://app.sandbox.midtrans.com/snap/v2/vtweb/...", nil
}
