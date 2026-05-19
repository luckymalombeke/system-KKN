package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"strings"
)

// MidtransSignature menghitung SHA512(order_id + status_code + gross_amount + serverKey).
func MidtransSignature(orderID, statusCode, grossAmount, serverKey string) string {
	input := orderID + statusCode + grossAmount + serverKey
	sum := sha512.Sum512([]byte(input))
	return hex.EncodeToString(sum[:])
}

// GrossAmountMatches membandingkan gross_amount notifikasi Midtrans dengan amount di DB.
func GrossAmountMatches(notificationAmount string, dbAmount int64) bool {
	normalized := strings.TrimSpace(notificationAmount)
	if normalized == "" {
		return false
	}

	parsed, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return false
	}

	return int64(parsed) == dbAmount
}
