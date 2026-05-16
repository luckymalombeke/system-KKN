# KKN System Management API

Sistem manajemen administrasi KKN (Kuliah Kerja Nyata) berbasis Go (Golang) dengan arsitektur yang bersih dan scalable.

## Tech Stack
- **Backend**: Go (Golang)
- **Framework**: Gin Gonic
- **ORM**: GORM
- **Database**: Supabase (PostgreSQL)
- **Auth**: JWT (JSON Web Token)
- **Payment**: Midtrans Integration (Planned)

## Folder Structure
- `cmd/`: Entry point aplikasi.
- `config/`: Konfigurasi environment dan sistem.
- `database/`: Koneksi database.
- `handlers/`: Controller untuk menangani HTTP request.
- `services/`: Logika bisnis utama.
- `repositories/`: Operasi langsung ke database.
- `models/`: Definisi data (Entity & DTO).

## Installation
1. Clone repository
2. Jalankan `go mod tidy`
3. Konfigurasi `.env`
4. Jalankan `go run cmd/api/main.go`
