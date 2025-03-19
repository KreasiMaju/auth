# Kreasimaju Auth

[![Go](https://github.com/kreasimaju/auth/actions/workflows/go.yml/badge.svg)](https://github.com/kreasimaju/auth/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kreasimaju/auth)](https://goreportcard.com/report/github.com/kreasimaju/auth)
[![Version](https://img.shields.io/github/v/tag/kreasimaju/auth)](https://github.com/kreasimaju/auth/releases)

Package autentikasi Go untuk aplikasi modern dengan dukungan lokal, OTP, dan OAuth.

## Fitur

- 🔒 Autentikasi lokal (email/password)
- 📱 Autentikasi dengan OTP (email, SMS, WhatsApp)
- 🌐 Autentikasi OAuth (Google, GitHub, dll.)
- 🔑 Manajemen token JWT
- 📞 Dukungan nomor telepon internasional
- 🔄 Manajemen reset password
- 🛡️ Validasi data robust

## Instalasi

```bash
go get github.com/kreasimaju/auth
```

## Konfigurasi

Buat file konfigurasi (JSON, YAML, atau ENV) dan muat saat inisialisasi:

```go
import "github.com/kreasimaju/auth"

func main() {
	// Muat konfigurasi dari file
	config, err := auth.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	// Inisialisasi autentikasi
	err = auth.Init(config)
	if err != nil {
		panic(err)
	}

	// Siap digunakan!
}
```

Contoh konfigurasi JSON:

```json
{
  "database": {
    "type": "mysql",
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "password",
    "database": "auth_db",
    "auto_migrate": true
  },
  "providers": {
    "google": {
      "enabled": true,
      "client_id": "your-client-id",
      "client_secret": "your-client-secret",
      "callback_url": "http://localhost:8080/auth/google/callback",
      "scopes": ["email", "profile"]
    },
    "local": true,
    "otp_auth": true
  },
  "jwt": {
    "secret": "your-jwt-secret",
    "expires_in": 86400
  },
  "otp": {
    "enabled": true,
    "default_type": "email",
    "length": 6,
    "expires_in": 300
  }
}
```

## Penggunaan Dasar

### Autentikasi Lokal

```go
// Registrasi pengguna baru
user, err := auth.RegisterLocal("user@example.com", "password123", "John", "Doe", "081234567890", "ID")

// Login
user, err := auth.Login("user@example.com", "password123")

// Generate JWT token
token, err := auth.GenerateToken(user)
```

### Autentikasi OTP

```go
// Meminta OTP untuk login
otpCode, err := auth.RequestOTPLogin("user@example.com", "email", "ID")

// Memverifikasi OTP
user, err := auth.VerifyOTPLogin("user@example.com", "email", "123456", "ID")
```

### Autentikasi OAuth

```go
// Redirect URL untuk autentikasi Google
url := auth.GetGoogleAuthURL("state")

// Proses callback dari Google
user, err := auth.HandleGoogleCallback(code)
```

## API Web

Package ini menyediakan handler HTTP siap pakai untuk Echo framework:

```go
import (
	"github.com/kreasimaju/auth"
	"github.com/labstack/echo/v4"
)

func main() {
	// Inisialisasi Echo
	e := echo.New()

	// Daftarkan route autentikasi
	auth.RegisterRoutes(e.Group("/auth"))

	// Jalankan server
	e.Start(":8080")
}
```

Ini akan menyediakan endpoint berikut:
- `POST /auth/register` - Registrasi pengguna
- `POST /auth/login` - Login pengguna
- `POST /auth/otp/request` - Request OTP
- `POST /auth/otp/verify` - Verifikasi OTP
- `GET /auth/oauth/google` - Login dengan Google
- Dan lainnya...

Untuk dokumentasi API lengkap, lihat [API Documentation](docs/API.md).

## Integrasi Provider SMS

Package ini mendukung integrasi dengan berbagai provider SMS, seperti:
- Twilio
- Vonage (Nexmo)
- Zenziva (Indonesia)
- Infobip
- Dan lainnya

Untuk panduan integrasi lebih lanjut, lihat [Panduan Integrasi SMS](docs/SMS_INTEGRATION.md).

## Build dari Source

Anda dapat membangun package ini langsung dari source:

```bash
# Clone repository
git clone https://github.com/kreasimaju/auth.git
cd auth

# Setup dependensi
make setup

# Jalankan pengujian
make test

# Build
make build
```

## Pengembangan

Kontribusi sangat dipersilakan! Untuk memulai pengembangan:

```bash
# Fork repository dan clone
git clone https://github.com/YOUR-USERNAME/auth.git
cd auth

# Setup
make setup

# Lakukan perubahan yang diinginkan...

# Format kode
make fmt

# Jalankan linter
make lint

# Jalankan pengujian
make test

# Buat pull request
```

## Versi

Lihat riwayat [Rilis](https://github.com/kreasimaju/auth/releases) untuk daftar perubahan dan versi.

## Lisensi

Package ini dilisensikan di bawah [MIT License](LICENSE). 