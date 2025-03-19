package models

import (
	"time"

	"gorm.io/gorm"
)

// OTPCode model untuk kode OTP
type OTPCode struct {
	gorm.Model
	UserID    uint       `gorm:"index" json:"user_id"`
	Code      string     `gorm:"type:varchar(10)" json:"code"`
	Type      string     `gorm:"type:varchar(20)" json:"type"`    // "sms", "email", "whatsapp"
	Target    string     `gorm:"type:varchar(255)" json:"target"` // email atau nomor telepon
	Purpose   string     `gorm:"type:varchar(50)" json:"purpose"` // "login", "register", "reset_password", "verify"
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	Attempts  int        `gorm:"default:0" json:"attempts"` // jumlah percobaan
	Valid     bool       `gorm:"default:true" json:"valid"`
}

// IsValid memeriksa apakah kode OTP masih valid
func (o *OTPCode) IsValid() bool {
	// Periksa apakah kode sudah digunakan
	if o.UsedAt != nil {
		return false
	}

	// Periksa apakah kode sudah kedaluwarsa
	if time.Now().After(o.ExpiresAt) {
		return false
	}

	// Periksa apakah kode masih valid
	return o.Valid
}

// MarkAsUsed menandai kode OTP sebagai sudah digunakan
func (o *OTPCode) MarkAsUsed(db *gorm.DB) error {
	now := time.Now()
	o.UsedAt = &now
	o.Valid = false
	return db.Save(o).Error
}

// IncrementAttempts menambah jumlah percobaan dan menginvalidasi jika melebihi batas
func (o *OTPCode) IncrementAttempts(db *gorm.DB, maxAttempts int) error {
	o.Attempts++
	if o.Attempts >= maxAttempts {
		o.Valid = false
	}
	return db.Save(o).Error
}
