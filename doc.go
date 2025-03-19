/*
Package auth menyediakan solusi autentikasi lengkap untuk aplikasi Go.

Package ini mendukung berbagai metode autentikasi termasuk lokal (email/password),
OTP (email, SMS, WhatsApp), dan OAuth (Google, GitHub, dll).

Fitur Utama:
  - Autentikasi lokal dengan email dan password
  - Autentikasi dengan OTP melalui email, SMS, atau WhatsApp
  - Autentikasi OAuth (Google, GitHub, dll.)
  - Dukungan JWT untuk otorisasi
  - Format nomor telepon internasional
  - Manajemen reset password
  - Middleware untuk perlindungan rute

Inisialisasi:

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

Autentikasi Lokal:

	// Registrasi pengguna baru
	user, err := auth.RegisterLocal("user@example.com", "password123", "John", "Doe", "081234567890", "ID")

	// Login
	user, err := auth.Login("user@example.com", "password123")

Autentikasi OTP:

	// Meminta OTP untuk login
	otpCode, err := auth.RequestOTPLogin("user@example.com", "email", "ID")

	// Memverifikasi OTP
	user, err := auth.VerifyOTPLogin("user@example.com", "email", "123456", "ID")

Autentikasi OAuth:

	// Redirect URL untuk autentikasi Google
	url := auth.GetGoogleAuthURL("state")

	// Proses callback dari Google
	user, err := auth.HandleGoogleCallback(code)

Integrasi dengan Web Framework (Echo):

	e := echo.New()

	// Daftarkan route autentikasi
	auth.RegisterRoutes(e.Group("/auth"))

	e.Start(":8080")

Untuk dokumentasi lebih lanjut kunjungi repository:
https://github.com/kreasimaju/auth
*/
package auth
