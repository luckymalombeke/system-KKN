package utils

import "testing"

func TestValidateJWTSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{"empty", "", true},
		{"default secret", "secret", true},
		{"too short", "abcd", true},
		{"strong", "kkn_dev_local_secret_min_32_chars_ok!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJWTSecret(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("secret=%q err=%v wantErr=%v", tt.secret, err, tt.wantErr)
			}
		})
	}
}
