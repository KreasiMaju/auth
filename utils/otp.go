package utils

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/kreasimaju/auth/config"
	"github.com/kreasimaju/auth/models"
	"gorm.io/gorm"
)

const (
	otpCharset = "0123456789"
)

// GenerateOTP menghasilkan kode OTP dengan panjang tertentu
func GenerateOTP(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	for i, b := range bytes {
		bytes[i] = otpCharset[b%byte(len(otpCharset))]
	}
	return string(bytes)
}

// CreateOTP membuat kode OTP baru dan menyimpannya di database
func CreateOTP(db *gorm.DB, userID uint, otpType, target, purpose string, cfg config.OTP) (*models.OTPCode, error) {
	// Hapus kode OTP yang sudah ada untuk tujuan yang sama
	db.Where("user_id = ? AND purpose = ? AND type = ? AND target = ? AND valid = ?", userID, purpose, otpType, target, true).
		Delete(&models.OTPCode{})

	// Buat kode OTP baru
	code := GenerateOTP(cfg.Length)
	expiresAt := time.Now().Add(time.Duration(cfg.ExpiresIn) * time.Second)

	otp := models.OTPCode{
		UserID:    userID,
		Code:      code,
		Type:      otpType,
		Target:    target,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		Valid:     true,
	}

	err := db.Create(&otp).Error
	if err != nil {
		return nil, err
	}

	return &otp, nil
}

// VerifyOTP memverifikasi kode OTP
func VerifyOTP(db *gorm.DB, userID uint, code, otpType, purpose string) (bool, error) {
	var otp models.OTPCode
	err := db.Where("user_id = ? AND code = ? AND type = ? AND purpose = ? AND valid = ?",
		userID, code, otpType, purpose, true).First(&otp).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	// Periksa validitas OTP
	if !otp.IsValid() {
		return false, nil
	}

	// Tandai sebagai digunakan
	if err := otp.MarkAsUsed(db); err != nil {
		return false, err
	}

	return true, nil
}

// SendOTP mengirim kode OTP melalui channel yang dipilih
func SendOTP(otpCode *models.OTPCode, cfg config.OTP) error {
	switch otpCode.Type {
	case "sms":
		return sendSMS(otpCode, cfg.SMS)
	case "email":
		return sendEmail(otpCode, cfg.Email)
	case "whatsapp":
		return sendWhatsApp(otpCode, cfg.WhatsApp)
	default:
		return fmt.Errorf("tipe OTP tidak dikenal: %s", otpCode.Type)
	}
}

// sendSMS mengirim kode OTP melalui SMS
func sendSMS(otpCode *models.OTPCode, cfg config.OTPProvider) error {
	// Implementasi pengiriman SMS
	// Gunakan provider SMS seperti Twilio, Nexmo, dll.
	return nil
}

// sendEmail mengirim kode OTP melalui email
func sendEmail(otpCode *models.OTPCode, cfg config.OTPProvider) error {
	// Implementasi pengiriman email
	// Gunakan SMTP atau layanan email seperti SendGrid, Mailgun, dll.
	return nil
}

// sendWhatsApp mengirim kode OTP melalui WhatsApp
func sendWhatsApp(otpCode *models.OTPCode, cfg config.OTPProvider) error {
	// Implementasi pengiriman WhatsApp
	// Gunakan API WhatsApp Business atau penyedia layanan seperti Twilio, dll.
	return nil
}
