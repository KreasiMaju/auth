package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/kreasimaju/auth/config"
	"github.com/kreasimaju/auth/middleware"
	"github.com/kreasimaju/auth/models"
	"github.com/kreasimaju/auth/providers"
	"github.com/kreasimaju/auth/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	configuration config.Config
)

// Init menginisialisasi package auth dengan konfigurasi yang diberikan
func Init(cfg config.Config) error {
	configuration = cfg

	// Set JWT secret
	utils.SetJWTSecret(cfg.JWT.Secret)

	// Menginisialisasi koneksi database
	_, err := utils.InitDB(cfg.Database)
	if err != nil {
		return err
	}

	// Inisialisasi provider auth
	initProviders(cfg.Providers)

	return nil
}

// initProviders mengatur provider auth berdasarkan konfigurasi
func initProviders(providers config.Providers) {
	// Inisialisasi provider OAuth
	if providers.Google.Enabled {
		initGoogleProvider(providers.Google)
	}

	if providers.Twitter.Enabled {
		initTwitterProvider(providers.Twitter)
	}

	if providers.GitHub.Enabled {
		initGitHubProvider(providers.GitHub)
	}

	if providers.Facebook.Enabled {
		initFacebookProvider(providers.Facebook)
	}
}

// initGoogleProvider menginisialisasi provider Google
func initGoogleProvider(config config.OAuth) {
	providers.InitGoogle(config)
}

// GoogleLoginURL mengembalikan URL untuk login Google
func GoogleLoginURL(state string) string {
	return providers.GoogleLoginURL(state)
}

// HandleGoogleCallback menangani callback dari Google OAuth
func HandleGoogleCallback(code string) (*models.User, error) {
	return providers.HandleGoogleCallback(code)
}

func initTwitterProvider(config config.OAuth) {
	// Implementasi inisialisasi Twitter OAuth provider
}

func initGitHubProvider(config config.OAuth) {
	// Implementasi inisialisasi GitHub OAuth provider
}

func initFacebookProvider(config config.OAuth) {
	// Implementasi inisialisasi Facebook OAuth provider
}

// GetConfig mengembalikan konfigurasi saat ini
func GetConfig() config.Config {
	return configuration
}

// Middleware functions

// Middleware mengembalikan handler middleware otentikasi untuk Echo
func Middleware() echo.MiddlewareFunc {
	return middleware.EchoAuthMiddleware()
}

// RoleMiddleware mengembalikan handler middleware peran untuk Echo
func RoleMiddleware(roles ...string) echo.MiddlewareFunc {
	return middleware.EchoRoleMiddleware(roles...)
}

// GinMiddleware mengembalikan handler middleware otentikasi untuk Gin
func GinMiddleware() gin.HandlerFunc {
	return middleware.GinAuthMiddleware()
}

// GinRoleMiddleware mengembalikan handler middleware peran untuk Gin
func GinRoleMiddleware(roles ...string) gin.HandlerFunc {
	return middleware.GinRoleMiddleware(roles...)
}

// FiberMiddleware mengembalikan handler middleware otentikasi untuk Fiber
func FiberMiddleware() fiber.Handler {
	return middleware.FiberAuthMiddleware()
}

// FiberRoleMiddleware mengembalikan handler middleware peran untuk Fiber
func FiberRoleMiddleware(roles ...string) fiber.Handler {
	return middleware.FiberRoleMiddleware(roles...)
}

// Route registration

// RegisterRoutes mendaftarkan rute otentikasi untuk Echo
func RegisterRoutes(e *echo.Echo) {
	// Registrasi rute autentikasi
	auth := e.Group("/auth")

	// Rute local auth
	auth.POST("/register", registerHandler)
	auth.POST("/login", loginHandler)
	auth.POST("/verify-email", verifyEmailHandler)
	auth.POST("/forgot-password", forgotPasswordHandler)
	auth.POST("/reset-password", resetPasswordHandler)

	// Rute OTP
	auth.POST("/request-otp", requestOTPHandler)
	auth.POST("/verify-otp", verifyOTPHandler)

	// Rute OAuth
	auth.GET("/google", googleAuthHandler)
	auth.GET("/google/callback", googleCallbackHandler)
	auth.GET("/twitter", twitterAuthHandler)
	auth.GET("/twitter/callback", twitterCallbackHandler)
	auth.GET("/github", githubAuthHandler)
	auth.GET("/github/callback", githubCallbackHandler)
	auth.GET("/facebook", facebookAuthHandler)
	auth.GET("/facebook/callback", facebookCallbackHandler)

	// Logout
	auth.POST("/logout", logoutHandler)
}

// ===== Autentikasi Lokal dan OTP =====

// RegisterLocal mendaftarkan pengguna baru dengan email dan password
func RegisterLocal(email, password, firstName, lastName, phone, defaultRegion string) (*models.User, error) {
	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Format nomor telepon ke format internasional jika ada
	formattedPhone := phone
	if phone != "" {
		formattedPhone, err = utils.FormatPhoneNumber(phone, defaultRegion)
		if err != nil {
			return nil, fmt.Errorf("format nomor telepon tidak valid: %v", err)
		}
	}

	// Buat pengguna baru
	user := models.User{
		Email:     email,
		Password:  hashedPassword,
		Phone:     formattedPhone,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user", // Default role
	}

	// Simpan ke database
	err = utils.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// LoginLocal melakukan autentikasi pengguna dengan email atau nomor telepon dan password
func LoginLocal(identifier, password string) (*models.User, error) {
	// Cari pengguna berdasarkan email atau telepon
	var user models.User

	// Coba cari dengan email
	err := utils.DB.Where("email = ?", identifier).First(&user).Error

	// Jika tidak ditemukan dengan email, coba cari dengan nomor telepon
	if err != nil && err == gorm.ErrRecordNotFound {
		err = utils.DB.Where("phone = ?", identifier).First(&user).Error
	}

	if err != nil {
		return nil, err
	}

	// Periksa password
	if !utils.CheckPassword(user.Password, password) {
		return nil, fmt.Errorf("password tidak valid")
	}

	// Update last login time
	now := time.Now()
	user.LastLogin = &now
	utils.DB.Save(&user)

	return &user, nil
}

// GenerateOTP membuat kode OTP baru
func GenerateOTP(userID uint, otpType, target, purpose string) (*models.OTPCode, error) {
	// Validasi tipe OTP
	if otpType != "sms" && otpType != "email" && otpType != "whatsapp" {
		return nil, fmt.Errorf("tipe OTP tidak valid")
	}

	// Cek apakah sudah ada OTP yang masih valid
	var existingOTP models.OTPCode
	err := utils.DB.Where("user_id = ? AND type = ? AND purpose = ? AND valid = ?",
		userID, otpType, purpose, true).First(&existingOTP).Error

	// Jika ada OTP yang masih aktif, gunakan kembali
	if err == nil && existingOTP.IsValid() {
		return &existingOTP, nil
	}

	// Tentukan panjang kode OTP
	codeLength := configuration.OTP.Length
	if codeLength <= 0 {
		codeLength = 6 // default 6 digit
	}

	// Generate kode OTP acak
	code := utils.GenerateOTP(codeLength)

	// Hitung waktu kedaluwarsa
	expiresAt := time.Now().Add(time.Duration(configuration.OTP.ExpiresIn) * time.Second)

	// Buat OTP baru
	otpCode := models.OTPCode{
		UserID:    userID,
		Code:      code,
		Type:      otpType,
		Target:    target,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		Valid:     true,
	}

	// Simpan ke database
	if err := utils.DB.Create(&otpCode).Error; err != nil {
		return nil, err
	}

	return &otpCode, nil
}

// VerifyOTP memverifikasi kode OTP
func VerifyOTP(userID uint, code, otpType, purpose string) (bool, error) {
	var otpCode models.OTPCode
	err := utils.DB.Where("user_id = ? AND code = ? AND type = ? AND purpose = ? AND valid = ?",
		userID, code, otpType, purpose, true).First(&otpCode).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("kode OTP tidak valid atau sudah kedaluwarsa")
		}
		return false, err
	}

	// Periksa validitas OTP
	if !otpCode.IsValid() {
		return false, fmt.Errorf("kode OTP tidak valid atau sudah kedaluwarsa")
	}

	// Tandai sebagai digunakan
	if err := otpCode.MarkAsUsed(utils.DB); err != nil {
		return false, err
	}

	return true, nil
}

// SendOTP mengirim kode OTP melalui channel yang dipilih
func SendOTP(otpCode *models.OTPCode) error {
	if otpCode == nil {
		return fmt.Errorf("OTP code is nil")
	}

	// Siapkan pesan OTP
	message := fmt.Sprintf("Kode OTP Anda adalah: %s. Kode berlaku selama %d detik.",
		otpCode.Code, configuration.OTP.ExpiresIn)

	// Implementasi pengiriman OTP
	switch otpCode.Type {
	case "email":
		// Implementasi pengiriman email
		// Untuk saat ini, tampilkan saja di log
		log.Printf("Sending OTP %s to email %s (Message: %s)", otpCode.Code, otpCode.Target, message)
	case "sms":
		// Implementasi pengiriman SMS
		// Untuk saat ini, tampilkan saja di log
		log.Printf("Sending OTP %s to SMS %s (Message: %s)", otpCode.Code, otpCode.Target, message)
	case "whatsapp":
		// Implementasi pengiriman WhatsApp
		// Untuk saat ini, tampilkan saja di log
		log.Printf("Sending OTP %s to WhatsApp %s (Message: %s)", otpCode.Code, otpCode.Target, message)
	default:
		return fmt.Errorf("tipe OTP tidak didukung")
	}

	return nil
}

// FindUserByEmail mencari pengguna berdasarkan email
func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := utils.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByPhone mencari pengguna berdasarkan nomor telepon
func FindUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := utils.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RequestOTPLogin meminta OTP untuk login
func RequestOTPLogin(contact, otpType, defaultRegion string) (*models.OTPCode, error) {
	var user models.User
	var err error
	formattedContact := contact

	// Cari pengguna berdasarkan kontak (email atau telepon)
	if otpType == "email" {
		err = utils.DB.Where("email = ?", contact).First(&user).Error
	} else { // sms atau whatsapp
		// Format nomor telepon jika perlu
		if utils.IsPhoneNumber(contact) {
			formattedContact, err = utils.FormatPhoneNumber(contact, defaultRegion)
			if err != nil {
				return nil, fmt.Errorf("format nomor telepon tidak valid: %v", err)
			}
		}
		err = utils.DB.Where("phone = ?", formattedContact).First(&user).Error
	}

	if err != nil {
		return nil, err
	}

	// Buat OTP
	otpCode, err := GenerateOTP(user.ID, otpType, formattedContact, "login")
	if err != nil {
		return nil, err
	}

	// Kirim OTP
	err = SendOTP(otpCode)
	if err != nil {
		return nil, err
	}

	return otpCode, nil
}

// VerifyOTPLogin memverifikasi OTP untuk login
func VerifyOTPLogin(contact, otpType, code, defaultRegion string) (*models.User, error) {
	var user models.User
	var err error
	formattedContact := contact

	// Cari pengguna berdasarkan kontak (email atau telepon)
	if otpType == "email" {
		err = utils.DB.Where("email = ?", contact).First(&user).Error
	} else { // sms atau whatsapp
		// Format nomor telepon jika perlu
		if utils.IsPhoneNumber(contact) {
			formattedContact, err = utils.FormatPhoneNumber(contact, defaultRegion)
			if err != nil {
				return nil, fmt.Errorf("format nomor telepon tidak valid: %v", err)
			}
		}
		err = utils.DB.Where("phone = ?", formattedContact).First(&user).Error
	}

	if err != nil {
		return nil, err
	}

	// Verifikasi OTP
	valid, err := VerifyOTP(user.ID, code, otpType, "login")
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, fmt.Errorf("kode OTP tidak valid")
	}

	// Update last login time
	now := time.Now()
	user.LastLogin = &now
	utils.DB.Save(&user)

	return &user, nil
}

// ===== Implementasi Handler =====

// Handler Google OAuth
var googleAuthHandler = func(c echo.Context) error {
	// Generate state untuk keamanan
	state := "random-state" // Idealnya gunakan random string dan simpan di session

	// Redirect ke URL login Google
	url := GoogleLoginURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Handler Google OAuth callback
var googleCallbackHandler = func(c echo.Context) error {
	// Dapatkan code dari query parameter
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Code parameter is required",
		})
	}

	// Dapatkan user dari Google callback
	user, err := HandleGoogleCallback(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to authenticate with Google: " + err.Error(),
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(*user, configuration.JWT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token: " + err.Error(),
		})
	}

	// Return token
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

// Handler untuk registrasi lokal
var registerHandler = func(c echo.Context) error {
	// Parse request
	var req struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Phone         string `json:"phone"`
		DefaultRegion string `json:"default_region"` // Kode negara 2 huruf, misal: "ID", "US", dll
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}

	// Cek apakah email sudah terdaftar
	_, err := FindUserByEmail(req.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Email already registered",
		})
	}

	// Jika nomor telepon disediakan, validasi formatnya
	if req.Phone != "" {
		// Validasi format telepon
		if !utils.IsValidPhoneNumber(req.Phone, req.DefaultRegion) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid phone number format",
			})
		}

		// Format telepon untuk penyimpanan
		formattedPhone, err := utils.FormatPhoneNumber(req.Phone, req.DefaultRegion)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid phone number: " + err.Error(),
			})
		}

		// Cek apakah nomor telepon sudah terdaftar
		_, err = FindUserByPhone(formattedPhone)
		if err == nil {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Phone number already registered",
			})
		}
	}

	// Registrasi pengguna
	user, err := RegisterLocal(req.Email, req.Password, req.FirstName, req.LastName, req.Phone, req.DefaultRegion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to register user: " + err.Error(),
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(*user, configuration.JWT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token: " + err.Error(),
		})
	}

	// Return token
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"phone":      user.Phone,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

// Handler untuk login lokal
var loginHandler = func(c echo.Context) error {
	// Parse request
	var req struct {
		Identifier string `json:"identifier"` // Email atau nomor telepon
		Password   string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validasi input
	if req.Identifier == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Identifier (email or phone) and password are required",
		})
	}

	// Login pengguna
	user, err := LoginLocal(req.Identifier, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid identifier or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(*user, configuration.JWT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token: " + err.Error(),
		})
	}

	// Return token
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"phone":      user.Phone,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

var (
	verifyEmailHandler      = func(c echo.Context) error { return nil }
	forgotPasswordHandler   = func(c echo.Context) error { return nil }
	resetPasswordHandler    = func(c echo.Context) error { return nil }
	twitterAuthHandler      = func(c echo.Context) error { return nil }
	twitterCallbackHandler  = func(c echo.Context) error { return nil }
	githubAuthHandler       = func(c echo.Context) error { return nil }
	githubCallbackHandler   = func(c echo.Context) error { return nil }
	facebookAuthHandler     = func(c echo.Context) error { return nil }
	facebookCallbackHandler = func(c echo.Context) error { return nil }
	logoutHandler           = func(c echo.Context) error { return nil }

	// Handler untuk request OTP
	requestOTPHandler = func(c echo.Context) error {
		// Parse request
		var req struct {
			Contact       string `json:"contact"`        // email atau nomor telepon
			Type          string `json:"type"`           // "sms", "email", "whatsapp"
			Purpose       string `json:"purpose"`        // "login", "register", "reset_password"
			DefaultRegion string `json:"default_region"` // Kode negara 2 huruf, default "ID"
		}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
		}

		// Validasi input
		if req.Contact == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Contact is required",
			})
		}

		if req.Type == "" {
			// Gunakan tipe default dari konfigurasi
			req.Type = configuration.OTP.DefaultType
		}

		if req.Type != "sms" && req.Type != "email" && req.Type != "whatsapp" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid OTP type",
			})
		}

		if req.Purpose == "" {
			req.Purpose = "login"
		}

		// Gunakan default region Indonesia jika tidak disediakan
		if req.DefaultRegion == "" {
			req.DefaultRegion = "ID"
		}

		// Jika untuk login, cek apakah pengguna sudah terdaftar
		if req.Purpose == "login" {
			// Request OTP untuk login
			_, err := RequestOTPLogin(req.Contact, req.Type, req.DefaultRegion)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.JSON(http.StatusNotFound, map[string]string{
						"error": "User not found",
					})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to request OTP: " + err.Error(),
				})
			}

			// Berhasil mengirim OTP
			return c.JSON(http.StatusOK, map[string]string{
				"message": "OTP has been sent",
			})
		} else if req.Purpose == "register" {
			// Untuk registrasi, cek apakah email/nomor sudah terdaftar
			var exists bool
			if req.Type == "email" {
				_, err := FindUserByEmail(req.Contact)
				exists = err == nil
			} else {
				// Format nomor telepon
				formattedPhone, err := utils.FormatPhoneNumber(req.Contact, req.DefaultRegion)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{
						"error": "Invalid phone number: " + err.Error(),
					})
				}
				_, err = FindUserByPhone(formattedPhone)
				exists = err == nil
			}

			if exists {
				return c.JSON(http.StatusConflict, map[string]string{
					"error": "Contact already registered",
				})
			}

			// Buat pengguna sementara atau gunakan session untuk menyimpan data OTP
			// Implementasi bergantung pada kebutuhan spesifik
			// ...

			return c.JSON(http.StatusOK, map[string]string{
				"message": "OTP has been sent for registration",
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid purpose",
		})
	}

	// Handler untuk verifikasi OTP
	verifyOTPHandler = func(c echo.Context) error {
		// Parse request
		var req struct {
			Contact       string `json:"contact"`        // email atau nomor telepon
			Type          string `json:"type"`           // "sms", "email", "whatsapp"
			Code          string `json:"code"`           // kode OTP
			Purpose       string `json:"purpose"`        // "login", "register", "reset_password"
			DefaultRegion string `json:"default_region"` // Kode negara 2 huruf, default "ID"
		}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
		}

		// Validasi input
		if req.Contact == "" || req.Code == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Contact and code are required",
			})
		}

		if req.Type == "" {
			req.Type = configuration.OTP.DefaultType
		}

		if req.Purpose == "" {
			req.Purpose = "login"
		}

		// Gunakan default region Indonesia jika tidak disediakan
		if req.DefaultRegion == "" {
			req.DefaultRegion = "ID"
		}

		// Verifikasi OTP untuk login
		if req.Purpose == "login" {
			user, err := VerifyOTPLogin(req.Contact, req.Type, req.Code, req.DefaultRegion)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid OTP: " + err.Error(),
				})
			}

			// Generate JWT token
			token, err := utils.GenerateJWT(*user, configuration.JWT)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to generate token: " + err.Error(),
				})
			}

			// Return token
			return c.JSON(http.StatusOK, map[string]interface{}{
				"token": token,
				"user": map[string]interface{}{
					"id":         user.ID,
					"email":      user.Email,
					"phone":      user.Phone,
					"first_name": user.FirstName,
					"last_name":  user.LastName,
				},
			})
		}

		// Implementasi verifikasi untuk registrasi dan reset password
		// ...

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid purpose",
		})
	}
)
