# KreasiMaju Auth

Package autentikasi yang dapat digunakan kembali untuk proyek-proyek berbasis Golang. Package ini menyediakan solusi otentikasi yang lengkap dengan dukungan berbagai provider dan migrasi database otomatis.

## Fitur

- 🔐 Autentikasi multi-provider (Google, Twitter, GitHub, dll)
- 📊 Migrasi database otomatis
- 🛡️ Middleware untuk perlindungan rute
- 🔄 Manajemen sesi dan token
- 📱 Dukungan untuk autentikasi berbasis JWT
- 👤 Manajemen profil pengguna
- 🔐 Reset password dan verifikasi email

## Instalasi

```bash
go get github.com/kreasimaju/auth
```

## Penggunaan Cepat

```go
package main

import (
	"github.com/kreasimaju/auth"
	"github.com/kreasimaju/auth/config"
)

func main() {
	// Inisialisasi auth dengan konfigurasi
	auth.Init(config.Config{
		Database: config.Database{
			Type:        "mysql", // atau "postgres", "sqlite", dll
			Host:        "localhost",
			Port:        3306,
			Username:    "root",
			Password:    "password",
			Database:    "myapp",
			AutoMigrate: true, // Jalankan migrasi otomatis
		},
		Providers: config.Providers{
			Google: config.OAuth{
				ClientID:     "your-client-id",
				ClientSecret: "your-client-secret",
				CallbackURL:  "http://localhost:8080/auth/google/callback",
			},
			Twitter: config.OAuth{
				// konfigurasi twitter
			},
			// provider lainnya...
		},
		OTP: config.OTP{
			Enabled:     true,
			DefaultType: "sms", // atau "email", "whatsapp"
			Length:      6,
			ExpiresIn:   300, // 5 menit
			SMS: config.OTPProvider{
				Enabled: true,
				APIKey:  "your-sms-api-key",
			},
			Email: config.OTPProvider{
				Enabled: true,
				APIKey:  "your-email-api-key",
			},
			WhatsApp: config.OTPProvider{
				Enabled: true,
				APIKey:  "your-whatsapp-api-key",
			},
		},
	})

	// Integrasi dengan framework Go web populer (contoh dengan Echo)
	e := echo.New()
	
	// Daftarkan middleware auth
	e.Use(auth.Middleware())
	
	// Tambahkan route autentikasi
	auth.RegisterRoutes(e)
	
	e.Start(":8080")
}
```

## Skema Database

Package ini akan membuat tabel-tabel berikut secara otomatis:

- `users` - Informasi dasar pengguna
- `user_providers` - Informasi provider autentikasi
- `sessions` - Sesi pengguna
- `tokens` - Token reset password, verifikasi email, dll

## Todo List

- [x] Setup struktur project dasar
- [x] Implementasi koneksi database dan migrasi otomatis menggunakan GORM
- [x] Membuat model User dan tabel terkait
- [x] Implementasi autentikasi lokal (email/password)
- [x] Implementasi JWT dan manajemen token
- [x] Integrasi provider OAuth (Google)
- [ ] Integrasi provider OAuth (Twitter)
- [ ] Integrasi provider OAuth (GitHub)
- [x] Implementasi middleware untuk framework Go populer (Echo, Gin, Fiber)
- [ ] Fitur reset password
- [x] Fitur OTP (SMS, Email, WhatsApp)
- [ ] Fitur verifikasi email
- [ ] Pengujian dan dokumentasi
- [ ] Publish ke GitHub
- [x] Contoh penggunaan dengan Echo
- [x] Contoh penggunaan dengan Gin
- [x] Contoh penggunaan dengan Fiber
- [ ] Implementasi hook dan event untuk kustomisasi

## Dokumentasi

Dokumentasi lengkap dapat ditemukan di [docs/README.md](docs/README.md)

## Lisensi

MIT 