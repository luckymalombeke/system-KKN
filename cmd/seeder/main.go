package main

import (
	"fmt"
	"kkn-system/config"
	"kkn-system/database"
	"kkn-system/models/entity"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("Memulai Seeder Admin...")

	// Load Config & DB
	cfg := config.LoadConfig()
	database.ConnectDB(cfg)

	// Pastikan tabel ada
	err := database.DB.AutoMigrate(&entity.Admin{})
	if err != nil {
		log.Fatal("Gagal migrasi:", err)
	}

	// Cek apakah admin sudah ada
	var count int64
	database.DB.Model(&entity.Admin{}).Count(&count)
	if count > 0 {
		fmt.Println("Admin sudah ada di database. Seeder dibatalkan.")
		return
	}

	// Generate Password Hash
	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal hash password:", err)
	}

	admin := entity.Admin{
		Name:     "Super Admin",
		Email:    "admin@kkn.ac.id",
		Password: string(hashedPassword),
		Role:     "superadmin",
	}

	err = database.DB.Create(&admin).Error
	if err != nil {
		log.Fatal("Gagal membuat admin:", err)
	}

	fmt.Println("Sukses membuat Super Admin!")
	fmt.Println("Email    : admin@kkn.ac.id")
	fmt.Println("Password : admin123")
}
