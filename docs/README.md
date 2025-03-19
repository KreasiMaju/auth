# Dokumentasi KreasiMaju Auth

Dokumentasi lengkap untuk package autentikasi KreasiMaju Auth.

## Daftar Isi

1. [Instalasi](#instalasi)
2. [Konfigurasi](#konfigurasi)
3. [Penggunaan Dasar](#penggunaan-dasar)
4. [Migrasi Database](#migrasi-database)
5. [Autentikasi JWT](#autentikasi-jwt)
6. [Provider OAuth](#provider-oauth)
7. [Middleware](#middleware)
8. [Contoh Kode](#contoh-kode)
9. [Autentikasi OTP](#autentikasi-otp)

## Instalasi

Instal package menggunakan Go modules:

```bash
go get github.com/kreasimaju/auth
```

## Konfigurasi

Package ini dikonfigurasi menggunakan struct `config.Config`. Berikut adalah contoh konfigurasi dasar:

```go
cfg := config.Config{
    Database: config.Database{
        Type:        "mysql",
        Host:        "localhost",
        Port:        3306,
        Username:    "root",
        Password:    "password",
        Database:    "myapp",
        AutoMigrate: true,
    },
    Providers: config.Providers{
        Google: config.OAuth{
            Enabled:      true,
            ClientID:     "your-client-id",
            ClientSecret: "your-client-secret",
            CallbackURL:  "http://localhost:8080/auth/google/callback",
        },
        Local: true, // Aktifkan autentikasi lokal (email/password)
    },
    JWT: config.JWT{
        Secret:    "your-jwt-secret-key",
        ExpiresIn: 86400, // 24 jam dalam detik
    },
}

err := auth.Init(cfg)
if err != nil {
    panic("Failed to initialize auth: " + err.Error())
}
```

## Penggunaan Dasar

### Inisialisasi

```go
import "github.com/kreasimaju/auth"

func main() {
    // Inisialisasi auth dengan konfigurasi
    err := auth.Init(cfg)
    if err != nil {
        panic("Failed to initialize auth: " + err.Error())
    }
    
    // ... kode aplikasi Anda ...
}
```

## Migrasi Database

Package ini mendukung migrasi database otomatis menggunakan GORM. Ketika `AutoMigrate` diatur ke `true` dalam konfigurasi database, tabel berikut akan dibuat secara otomatis:

- `users` - Tabel pengguna dasar
- `user_providers` - Informasi provider autentikasi
- `sessions` - Sesi pengguna
- `tokens` - Token untuk reset password, verifikasi email, dll

## Autentikasi JWT

Package ini menggunakan JSON Web Token (JWT) untuk otentikasi. Token dihasilkan saat pengguna login dan divalidasi di setiap permintaan yang dilindungi.

```go
// Menghasilkan token untuk pengguna
tokenString, err := utils.GenerateJWT(user, cfg.JWT)

// Memvalidasi token
token, err := utils.ValidateJWT(tokenString)
```

## Provider OAuth

### Google

```go
// Mendapatkan URL untuk login Google
url := auth.GoogleLoginURL("state-value")

// Menangani callback Google
user, err := auth.HandleGoogleCallback(code)
```

### Implementasi Provider Lainnya

Package ini juga mendukung provider lain seperti Twitter, GitHub, dan Facebook. Implementasi serupa dengan Google.

## Middleware

### Echo Framework

```go
import (
    "github.com/kreasimaju/auth"
    "github.com/labstack/echo/v4"
)

func setupRoutes() {
    e := echo.New()
    
    // Middleware otentikasi
    api := e.Group("/api")
    api.Use(auth.Middleware())
    
    // Middleware peran
    admin := api.Group("/admin")
    admin.Use(auth.RoleMiddleware("admin"))
}
```

### Gin Framework

```go
import (
    "github.com/kreasimaju/auth"
    "github.com/gin-gonic/gin"
)

func setupRoutes() {
    r := gin.Default()
    
    // Middleware otentikasi
    api := r.Group("/api")
    api.Use(auth.GinMiddleware())
    
    // Middleware peran
    admin := api.Group("/admin")
    admin.Use(auth.GinRoleMiddleware("admin"))
}
```

### Fiber Framework

```go
import (
    "github.com/kreasimaju/auth"
    "github.com/gofiber/fiber/v2"
)

func setupRoutes() {
    app := fiber.New()
    
    // Middleware otentikasi
    api := app.Group("/api")
    api.Use(auth.FiberMiddleware())
    
    // Middleware peran
    admin := api.Group("/admin")
    admin.Use(auth.FiberRoleMiddleware("admin"))
}
```

## Contoh Kode

Lihat direktori `examples/` untuk contoh implementasi lengkap dengan berbagai framework.

- [Contoh Echo](/examples/echo)
- [Contoh Gin](/examples/gin)
- [Contoh Fiber](/examples/fiber)

## Autentikasi OTP

Package ini mendukung autentikasi menggunakan One-Time Password (OTP) melalui berbagai saluran.

### Konfigurasi OTP

```go
cfg := config.Config{
    // ... konfigurasi lainnya ...
    
    OTP: config.OTP{
        Enabled:     true,
        DefaultType: "sms", // atau "email", "whatsapp"
        Length:      6, // panjang kode OTP
        ExpiresIn:   300, // 5 menit dalam detik
        SMS: config.OTPProvider{
            Enabled:  true,
            APIKey:   "your-sms-api-key",
            Sender:   "KreasiMaju",
            Template: "Kode OTP Anda adalah {code}. Kode ini berlaku selama 5 menit.",
        },
        Email: config.OTPProvider{
            Enabled:  true,
            APIKey:   "your-email-api-key",
            Sender:   "no-reply@example.com",
            Template: "Kode OTP Anda adalah {code}. Kode ini berlaku selama 5 menit.",
        },
        WhatsApp: config.OTPProvider{
            Enabled:  true,
            APIKey:   "your-whatsapp-api-key",
            Sender:   "your-whatsapp-number",
            Template: "Kode OTP Anda adalah {code}. Kode ini berlaku selama 5 menit.",
        },
    },
}
```

### Meminta OTP

```go
// Meminta OTP untuk login dengan email
otpCode, err := auth.RequestOTPLogin("user@example.com", "email", "")

// Meminta OTP untuk login dengan SMS (Indonesia)
otpCode, err := auth.RequestOTPLogin("+6281234567890", "sms", "ID")

// Meminta OTP untuk login dengan WhatsApp (dengan format nomor lokal dan kode negara)
otpCode, err := auth.RequestOTPLogin("081234567890", "whatsapp", "ID")
```

### Memverifikasi OTP

```go
// Memverifikasi OTP untuk login dengan email
user, err := auth.VerifyOTPLogin("user@example.com", "email", "123456", "")

// Memverifikasi OTP untuk login dengan SMS (Indonesia)
user, err := auth.VerifyOTPLogin("+6281234567890", "sms", "123456", "ID")

// Memverifikasi OTP untuk login dengan WhatsApp (dengan format nomor lokal dan kode negara)
user, err := auth.VerifyOTPLogin("081234567890", "whatsapp", "123456", "ID")
``` 