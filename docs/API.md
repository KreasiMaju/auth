# API Documentation for Kreasimaju Auth

## Base URL
```
https://your-api-domain.com/auth
```

## Autentikasi
Semua permintaan yang memerlukan autentikasi harus menyertakan header `Authorization` dengan token JWT.

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Endpoints

### Registrasi Pengguna

**Endpoint:** `POST /register`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "081234567890",
  "default_region": "ID"
}
```

- `email`: Alamat email pengguna (wajib)
- `password`: Kata sandi (minimal 8 karakter) (wajib)
- `first_name`: Nama depan (wajib)
- `last_name`: Nama belakang (opsional)
- `phone`: Nomor telepon (opsional)
- `default_region`: Kode negara 2 huruf untuk format nomor telepon (default: "ID")

**Response Sukses (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "phone": "+6281234567890",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Response Error (409 Conflict):**
```json
{
  "error": "Email already registered"
}
```

### Login Pengguna

**Endpoint:** `POST /login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response Sukses (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "phone": "+6281234567890",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Response Error (401 Unauthorized):**
```json
{
  "error": "Invalid email or password"
}
```

### Request OTP

**Endpoint:** `POST /otp/request`

**Request:**
```json
{
  "contact": "user@example.com",
  "type": "email",
  "purpose": "login",
  "default_region": "ID"
}
```

- `contact`: Email atau nomor telepon (wajib)
- `type`: Tipe OTP: "email", "sms", atau "whatsapp" (default dari konfigurasi)
- `purpose`: Tujuan OTP: "login", "register", atau "reset_password" (default: "login")
- `default_region`: Kode negara 2 huruf untuk format nomor telepon (default: "ID")

**Response Sukses (200 OK):**
```json
{
  "message": "OTP has been sent"
}
```

**Response Error (404 Not Found):**
```json
{
  "error": "User not found"
}
```

### Verifikasi OTP

**Endpoint:** `POST /otp/verify`

**Request:**
```json
{
  "contact": "user@example.com",
  "type": "email",
  "code": "123456",
  "purpose": "login",
  "default_region": "ID"
}
```

- `contact`: Email atau nomor telepon (wajib)
- `type`: Tipe OTP: "email", "sms", atau "whatsapp" (default dari konfigurasi)
- `code`: Kode OTP (wajib)
- `purpose`: Tujuan OTP: "login", "register", atau "reset_password" (default: "login")
- `default_region`: Kode negara 2 huruf untuk format nomor telepon (default: "ID")

**Response Sukses (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "phone": "+6281234567890",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Response Error (401 Unauthorized):**
```json
{
  "error": "Invalid OTP"
}
```

### Login dengan Google OAuth

**Endpoint:** `GET /oauth/google`

Mengarahkan pengguna ke halaman login Google.

**Callback URL:** `GET /oauth/google/callback`

**Response Sukses (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

### Permintaan Reset Password

**Endpoint:** `POST /password/reset/request`

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response Sukses (200 OK):**
```json
{
  "message": "Password reset link has been sent"
}
```

### Reset Password

**Endpoint:** `POST /password/reset`

**Request:**
```json
{
  "token": "reset_token_from_email",
  "password": "new_password123"
}
```

**Response Sukses (200 OK):**
```json
{
  "message": "Password has been reset"
}
```

**Response Error (400 Bad Request):**
```json
{
  "error": "Invalid or expired token"
}
```

### Validasi Token

**Endpoint:** `GET /validate`

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response Sukses (200 OK):**
```json
{
  "valid": true,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "phone": "+6281234567890",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Response Error (401 Unauthorized):**
```json
{
  "valid": false,
  "error": "Invalid or expired token"
}
```

## Format Nomor Telepon

Nomor telepon diformat ke standar internasional E.164. Parameter `default_region` (2 huruf kode negara) digunakan untuk menentukan format nomor telepon.

**Contoh:**
- Nomor Indonesia: "081234567890" dengan `default_region: "ID"` akan diformat sebagai "+6281234567890"
- Nomor AS: "2125551212" dengan `default_region: "US"` akan diformat sebagai "+12125551212"

## Kode Status HTTP

- `200 OK`: Permintaan berhasil
- `400 Bad Request`: Parameter tidak valid
- `401 Unauthorized`: Autentikasi gagal
- `404 Not Found`: Sumber daya tidak ditemukan
- `409 Conflict`: Konflik dengan sumber daya yang ada (misalnya, email sudah terdaftar)
- `500 Internal Server Error`: Kesalahan server 