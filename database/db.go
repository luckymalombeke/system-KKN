package database

import (
	"fmt"
	"kkn-system/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg config.Config) {
	dsn := fmt.Sprintf("host='%s' user='%s' password='%s' dbname='%s' port='%s' sslmode='%s'",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("Database connection established")
}
