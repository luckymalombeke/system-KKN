package utils

import (
	"testing"
)

func TestGenerateOTP6(t *testing.T) {
	for i := 0; i < 20; i++ {
		otp, err := GenerateOTP6()
		if err != nil {
			t.Fatal(err)
		}
		if len(otp) != 6 {
			t.Fatalf("expected length 6, got %q", otp)
		}
		for _, c := range otp {
			if c < '0' || c > '9' {
				t.Fatalf("non-digit in otp: %q", otp)
			}
		}
	}
}
