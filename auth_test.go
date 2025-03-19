package auth

import (
	"testing"
	"time"

	"github.com/kreasimaju/auth/models"
	"github.com/kreasimaju/auth/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB menyiapkan database untuk pengujian
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Migrasi skema untuk pengujian
	err = db.AutoMigrate(&models.User{}, &models.OTPCode{}, &models.PasswordReset{}, &models.OAuth{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Setel database global untuk pengujian
	utils.DB = db

	return db
}

// TestRegisterLocal menguji fungsi RegisterLocal
func TestRegisterLocal(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	// Test case 1: Registrasi pengguna berhasil
	t.Run("Successful Registration", func(t *testing.T) {
		user, err := RegisterLocal("test@example.com", "password123", "John", "Doe", "081234567890", "ID")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "+6281234567890", user.Phone) // Diformat ke standar internasional
	})

	// Test case 2: Email duplikat
	t.Run("Duplicate Email", func(t *testing.T) {
		_, err := RegisterLocal("test@example.com", "anotherpassword", "Jane", "Smith", "087654321", "ID")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email sudah terdaftar")
	})

	// Test case 3: Nomor telepon duplikat
	t.Run("Duplicate Phone", func(t *testing.T) {
		_, err := RegisterLocal("another@example.com", "password123", "Another", "User", "081234567890", "ID")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nomor telepon sudah terdaftar")
	})
}

// TestLogin menguji fungsi Login
func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	// Siapkan pengguna untuk pengujian
	_, err := RegisterLocal("login@example.com", "password123", "Login", "User", "081234567890", "ID")
	assert.NoError(t, err)

	// Test case 1: Login berhasil
	t.Run("Successful Login", func(t *testing.T) {
		user, err := Login("login@example.com", "password123")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "login@example.com", user.Email)
		assert.NotNil(t, user.LastLogin)
	})

	// Test case 2: Email tidak ditemukan
	t.Run("Email Not Found", func(t *testing.T) {
		_, err := Login("notfound@example.com", "password123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email atau password tidak valid")
	})

	// Test case 3: Password salah
	t.Run("Wrong Password", func(t *testing.T) {
		_, err := Login("login@example.com", "wrongpassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email atau password tidak valid")
	})
}

// TestOTPFlow menguji alur OTP untuk login
func TestOTPFlow(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	// Siapkan pengguna untuk pengujian
	user, err := RegisterLocal("otp@example.com", "password123", "OTP", "User", "081234567890", "ID")
	assert.NoError(t, err)

	// Test case 1: Request OTP berhasil
	t.Run("Request OTP", func(t *testing.T) {
		otpCode, err := RequestOTPLogin(user.Email, "email", "")
		assert.NoError(t, err)
		assert.NotNil(t, otpCode)
		assert.Equal(t, user.ID, otpCode.UserID)
		assert.Equal(t, "email", otpCode.Type)
		assert.Equal(t, "login", otpCode.Purpose)
		assert.True(t, otpCode.ExpiresAt.After(time.Now()))
	})

	// Test case 2: Verifikasi OTP berhasil
	t.Run("Verify OTP", func(t *testing.T) {
		// Request OTP terlebih dahulu
		otpCode, err := RequestOTPLogin(user.Email, "email", "")
		assert.NoError(t, err)

		// Verifikasi OTP
		verifiedUser, err := VerifyOTPLogin(user.Email, "email", otpCode.Code, "")
		assert.NoError(t, err)
		assert.NotNil(t, verifiedUser)
		assert.Equal(t, user.ID, verifiedUser.ID)
		assert.NotNil(t, verifiedUser.LastLogin)
	})

	// Test case 3: OTP tidak valid
	t.Run("Invalid OTP", func(t *testing.T) {
		_, err := VerifyOTPLogin(user.Email, "email", "123456", "")
		assert.Error(t, err)
	})
}

// TestPhoneValidation menguji validasi dan pemformatan nomor telepon
func TestPhoneValidation(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	testCases := []struct {
		name        string
		phone       string
		region      string
		expectError bool
		expected    string
	}{
		{"Valid Indonesian Phone", "081234567890", "ID", false, "+6281234567890"},
		{"Valid Indonesian Phone with +", "+6281234567890", "ID", false, "+6281234567890"},
		{"Valid Indonesian Phone with Different Format", "0812-3456-7890", "ID", false, "+6281234567890"},
		{"Valid US Phone", "2125551212", "US", false, "+12125551212"},
		{"Invalid Phone Format", "abcdef", "ID", true, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formattedPhone, err := utils.FormatPhoneNumber(tc.phone, tc.region)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, formattedPhone)
			}
		})
	}
}
