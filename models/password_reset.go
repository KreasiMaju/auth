package models

import (
	"time"

	"gorm.io/gorm"
)

// PasswordReset model untuk reset password
type PasswordReset struct {
	gorm.Model
	UserID    uint      `gorm:"index" json:"user_id"`
	Token     string    `gorm:"type:varchar(100);uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
}

// IsValid memeriksa apakah token reset masih valid
func (pr *PasswordReset) IsValid() bool {
	return !pr.Used && time.Now().Before(pr.ExpiresAt)
}

// MarkAsUsed menandai token sebagai sudah digunakan
func (pr *PasswordReset) MarkAsUsed(db *gorm.DB) error {
	pr.Used = true
	return db.Save(pr).Error
}
