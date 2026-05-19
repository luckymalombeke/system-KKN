package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	DevLog   bool
}

type Config struct {
	DBHost                      string
	DBUser                      string
	DBPassword                  string
	DBName                      string
	DBPort                      string
	DBSSLMode                   string
	JWTSecret                   string
	AppPort                     string
	MidtransServerKey           string
	MidtransIsProduction        bool
	MidtransSkipSignatureVerify bool
	PaymentExpiryHours          int
	SMTP                        SMTPConfig
	OTPExpiryMinutes            int
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return Config{
		DBHost:                      getEnv("DB_HOST", "localhost"),
		DBUser:                      getEnv("DB_USER", "postgres"),
		DBPassword:                  getEnv("DB_PASSWORD", ""),
		DBName:                      getEnv("DB_NAME", "kkn_db"),
		DBPort:                      getEnv("DB_PORT", "5432"),
		DBSSLMode:                   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:                   getEnv("JWT_SECRET", ""),
		AppPort:                     getEnv("PORT", "8081"),
		MidtransServerKey:           strings.Trim(strings.TrimSpace(getEnv("MIDTRANS_SERVER_KEY", "")), `"'`),
		MidtransIsProduction:        getEnvBool("MIDTRANS_IS_PRODUCTION", false),
		MidtransSkipSignatureVerify: getEnvBool("MIDTRANS_SKIP_SIGNATURE_VERIFY", false),
		PaymentExpiryHours:          getEnvInt("PAYMENT_EXPIRY_HOURS", 24),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnvInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
			DevLog:   getEnvBool("OTP_DEV_LOG", false),
		},
		OTPExpiryMinutes: getEnvInt("OTP_EXPIRY_MINUTES", 5),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
