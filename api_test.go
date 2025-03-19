package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kreasimaju/auth/models"
	"github.com/kreasimaju/auth/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupAPITest menyiapkan server Echo dan database untuk pengujian API
func setupAPITest(t *testing.T) (*echo.Echo, *gorm.DB) {
	// Siapkan database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Migrasi skema
	err = db.AutoMigrate(&models.User{}, &models.OTPCode{}, &models.PasswordReset{}, &models.OAuth{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Setel database global
	utils.DB = db

	// Siapkan server Echo
	e := echo.New()

	// Daftarkan route API
	e.POST("/auth/register", registerHandler)
	e.POST("/auth/login", loginHandler)
	e.POST("/auth/otp/request", requestOTPHandler)
	e.POST("/auth/otp/verify", verifyOTPHandler)

	return e, db
}

// TestRegisterAPI menguji endpoint API registrasi
func TestRegisterAPI(t *testing.T) {
	e, db := setupAPITest(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	// Test case: Registrasi berhasil
	t.Run("Successful Registration", func(t *testing.T) {
		// Buat request
		reqBody := map[string]interface{}{
			"email":          "apitest@example.com",
			"password":       "password123",
			"first_name":     "API",
			"last_name":      "Test",
			"phone":          "081234567890",
			"default_region": "ID",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Panggil handler
		err := registerHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Periksa respons
		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "token")
		assert.Contains(t, resp, "user")
	})

	// Test case: Email duplikat
	t.Run("Duplicate Email", func(t *testing.T) {
		// Buat request dengan email yang sama
		reqBody := map[string]interface{}{
			"email":          "apitest@example.com",
			"password":       "password123",
			"first_name":     "Another",
			"last_name":      "User",
			"phone":          "087654321",
			"default_region": "ID",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Panggil handler
		err := registerHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		// Periksa respons
		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "error")
	})
}

// TestLoginAPI menguji endpoint API login
func TestLoginAPI(t *testing.T) {
	e, db := setupAPITest(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	// Siapkan pengguna untuk pengujian
	_, err := RegisterLocal("logintest@example.com", "password123", "Login", "Test", "081234567890", "ID")
	assert.NoError(t, err)

	// Test case: Login berhasil
	t.Run("Successful Login", func(t *testing.T) {
		// Buat request
		reqBody := map[string]interface{}{
			"identifier": "logintest@example.com",
			"password":   "password123",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Panggil handler
		err := loginHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Periksa respons
		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "token")
		assert.Contains(t, resp, "user")
	})

	// Test case: Kredensial salah
	t.Run("Invalid Credentials", func(t *testing.T) {
		// Buat request dengan password salah
		reqBody := map[string]interface{}{
			"identifier": "logintest@example.com",
			"password":   "wrongpassword",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Panggil handler
		err := loginHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Periksa respons
		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "error")
	})
}

// TestOTPAPI menguji endpoint API OTP
func TestOTPAPI(t *testing.T) {
	e, db := setupAPITest(t)
	defer func() {
		dbSQL, _ := db.DB()
		dbSQL.Close()
	}()

	// Siapkan pengguna untuk pengujian
	_, err := RegisterLocal("otptest@example.com", "password123", "OTP", "Test", "081234567890", "ID")
	assert.NoError(t, err)

	var otpCode string

	// Test case: Request OTP berhasil
	t.Run("Request OTP", func(t *testing.T) {
		// Buat request
		reqBody := map[string]interface{}{
			"contact":        "otptest@example.com",
			"type":           "email",
			"purpose":        "login",
			"default_region": "ID",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/otp/request", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Panggil handler
		err := requestOTPHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Ambil kode OTP dari database untuk verifikasi selanjutnya
		var otpRecord models.OTPCode
		err = db.Where("type = ? AND target = ? AND purpose = ?", "email", "otptest@example.com", "login").First(&otpRecord).Error
		assert.NoError(t, err)
		otpCode = otpRecord.Code
	})

	// Test case: Verifikasi OTP berhasil
	t.Run("Verify OTP", func(t *testing.T) {
		// Buat request
		reqBody := map[string]interface{}{
			"contact":        "otptest@example.com",
			"type":           "email",
			"code":           otpCode,
			"purpose":        "login",
			"default_region": "ID",
		}
		jsonData, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/otp/verify", bytes.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Panggil handler
		err := verifyOTPHandler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Periksa respons
		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp, "token")
		assert.Contains(t, resp, "user")
	})
}
