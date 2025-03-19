package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Email         string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password      string         `gorm:"type:varchar(255)" json:"-"`
	Phone         string         `gorm:"type:varchar(20);index" json:"phone"`
	FirstName     string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName      string         `gorm:"type:varchar(100)" json:"last_name"`
	Role          string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsVerified    bool           `gorm:"default:false" json:"is_verified"`
	LastLogin     *time.Time     `json:"last_login"`
	Providers     []UserProvider `json:"providers"`
	Sessions      []Session      `json:"-"`
	PasswordReset []Token        `gorm:"polymorphic:Owner;polymorphicValue:password_reset" json:"-"`
	EmailVerify   []Token        `gorm:"polymorphic:Owner;polymorphicValue:email_verify" json:"-"`
	OTPCodes      []OTPCode      `json:"-"`
}

// UserProvider model untuk provider autentikasi
type UserProvider struct {
	gorm.Model
	UserID       uint       `gorm:"index" json:"user_id"`
	ProviderName string     `gorm:"type:varchar(50)" json:"provider_name"`
	ProviderID   string     `gorm:"type:varchar(100)" json:"provider_id"`
	AccessToken  string     `gorm:"type:text" json:"-"`
	RefreshToken string     `gorm:"type:text" json:"-"`
	ExpiresAt    *time.Time `json:"-"`
	Data         string     `gorm:"type:text" json:"-"` // JSON data dari provider
}

// Session model untuk sesi pengguna
type Session struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex" json:"token"`
	UserAgent string    `gorm:"type:varchar(255)" json:"user_agent"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Token model untuk reset password, verifikasi email, dll
type Token struct {
	gorm.Model
	OwnerID   uint       `gorm:"index" json:"owner_id"`
	OwnerType string     `gorm:"type:varchar(50)" json:"owner_type"`
	Token     string     `gorm:"type:varchar(255);uniqueIndex" json:"token"`
	Type      string     `gorm:"type:varchar(50)" json:"type"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
}
