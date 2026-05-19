package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kkn-system/config"
	"kkn-system/utils"
	"net/http"
	"strings"
	"time"
)

// MidtransNotification payload HTTP notification Midtrans.
type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	TransactionStatus string `json:"transaction_status"`
	SignatureKey      string `json:"signature_key"`
	FraudStatus       string `json:"fraud_status"`
}

type SnapCustomerDetail struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

type PaymentService interface {
	CreateSnapTransaction(orderID string, amount int64, customer SnapCustomerDetail, expiredAt time.Time) (redirectURL string, err error)
	VerifyNotification(notification MidtransNotification) error
}

type paymentService struct {
	serverKey           string
	isProduction        bool
	skipSignatureVerify bool
	httpClient          *http.Client
}

func NewPaymentService(cfg config.Config) PaymentService {
	return &paymentService{
		serverKey:           cfg.MidtransServerKey,
		isProduction:        cfg.MidtransIsProduction,
		skipSignatureVerify: cfg.MidtransSkipSignatureVerify,
		httpClient:          &http.Client{},
	}
}

func (s *paymentService) CreateSnapTransaction(orderID string, amount int64, customer SnapCustomerDetail, expiredAt time.Time) (string, error) {
	if s.serverKey == "" {
		return "", errors.New("MIDTRANS_SERVER_KEY belum dikonfigurasi")
	}
	if err := validateMidtransServerKey(s.serverKey, s.isProduction); err != nil {
		return "", err
	}
	if amount <= 0 {
		return "", errors.New("nominal pembayaran tidak valid")
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60)
	}
	start := time.Now().In(loc)
	durationHours := int(expiredAt.In(loc).Sub(start).Hours())
	if durationHours < 1 {
		durationHours = 1
	}

	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": amount,
		},
		"customer_details": map[string]string{
			"first_name": customer.FirstName,
			"email":      customer.Email,
		},
		"expiry": map[string]interface{}{
			"start_time": start.Format("2006-01-02 15:04:05 -0700"),
			"unit":       "hour",
			"duration":   durationHours,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, s.snapBaseURL(), bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s.serverKey+":")))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal menghubungi Midtrans: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", errors.New("Midtrans menolak Server Key (401). Salin Server Key persis dari dashboard (Mid-server-..., bukan Client Key), tanpa menambah SB-, lalu restart backend")
		}
		return "", fmt.Errorf("Midtrans menolak transaksi (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var snapResp struct {
		RedirectURL string `json:"redirect_url"`
	}
	if err := json.Unmarshal(respBody, &snapResp); err != nil {
		return "", err
	}
	if snapResp.RedirectURL == "" {
		return "", errors.New("Midtrans tidak mengembalikan redirect_url")
	}

	return snapResp.RedirectURL, nil
}

func (s *paymentService) VerifyNotification(notification MidtransNotification) error {
	if s.skipSignatureVerify {
		return nil
	}

	if s.serverKey == "" {
		return errors.New("MIDTRANS_SERVER_KEY belum dikonfigurasi")
	}
	if notification.OrderID == "" || notification.StatusCode == "" || notification.GrossAmount == "" {
		return errors.New("payload notifikasi tidak lengkap")
	}
	if notification.SignatureKey == "" {
		return errors.New("signature_key wajib ada")
	}

	expected := utils.MidtransSignature(
		notification.OrderID,
		notification.StatusCode,
		notification.GrossAmount,
		s.serverKey,
	)

	if !strings.EqualFold(expected, notification.SignatureKey) {
		return errors.New("signature_key tidak valid")
	}

	return nil
}

func validateMidtransServerKey(key string, _ bool) error {
	lower := strings.ToLower(key)
	if strings.Contains(lower, "mid-client") {
		return errors.New("yang terpasang terlihat Client Key; gunakan Server Key (Mid-server-...) dari dashboard")
	}
	if strings.HasPrefix(key, "Mid-server-") || strings.HasPrefix(key, "SB-Mid-server-") {
		return nil
	}
	return errors.New("Server Key tidak valid; salin Server Key dari dashboard (Mid-server-... atau SB-Mid-server-...)")
}

func (s *paymentService) snapBaseURL() string {
	if s.isProduction {
		return "https://app.midtrans.com/snap/v1/transactions"
	}
	return "https://app.sandbox.midtrans.com/snap/v1/transactions"
}

// MapMidtransTransactionStatus mengubah status Midtrans ke status aplikasi.
func MapMidtransTransactionStatus(notification MidtransNotification) string {
	switch notification.TransactionStatus {
	case "settlement", "capture":
		if notification.StatusCode == "200" {
			fraud := strings.ToLower(strings.TrimSpace(notification.FraudStatus))
			if fraud == "" || fraud == "accept" {
				return "success"
			}
		}
		return "pending"
	case "expire", "cancel", "deny", "failure":
		return "failed"
	default:
		return "pending"
	}
}
