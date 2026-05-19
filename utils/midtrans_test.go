package utils

import "testing"

func TestMidtransSignature(t *testing.T) {
	orderID := "INV-TEST-001"
	statusCode := "200"
	grossAmount := "500000.00"
	serverKey := "SB-Mid-server-test-key"

	sig := MidtransSignature(orderID, statusCode, grossAmount, serverKey)
	if len(sig) != 128 {
		t.Fatalf("expected SHA512 hex length 128, got %d", len(sig))
	}

	if MidtransSignature(orderID, statusCode, grossAmount, serverKey) != sig {
		t.Fatal("signature should be deterministic")
	}
}

func TestGrossAmountMatches(t *testing.T) {
	if !GrossAmountMatches("500000.00", 500000) {
		t.Fatal("500000.00 should match 500000")
	}
	if GrossAmountMatches("499999.00", 500000) {
		t.Fatal("499999.00 should not match 500000")
	}
}
