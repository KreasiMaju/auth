# Panduan Integrasi SMS API

Dokumen ini menjelaskan cara mengintegrasikan berbagai provider SMS API dengan package `kreasimaju-auth`.

## Provider yang Didukung

Berikut adalah beberapa provider SMS yang populer dan dapat diintegrasikan:

1. Twilio
2. Vonage (sebelumnya Nexmo)
3. Zenziva (Indonesia)
4. Wavecell (sekarang 8x8)
5. Infobip

## Struktur Integrasi

Package `kreasimaju-auth` menggunakan interface untuk provider SMS, sehingga Anda dapat mengimplementasikan provider pilihan Anda:

```go
// SMSProvider adalah interface untuk provider SMS
type SMSProvider interface {
    SendSMS(to, message string) error
}
```

## Integrasi dengan Twilio

### Langkah 1: Buat Akun Twilio
Daftar di [Twilio](https://www.twilio.com) dan dapatkan:
- Account SID
- Auth Token
- Nomor Telepon Twilio

### Langkah 2: Instal Modul Twilio Go
```bash
go get github.com/twilio/twilio-go
```

### Langkah 3: Implementasi Provider

Buat file `providers/twilio.go`:

```go
package providers

import (
    "github.com/twilio/twilio-go"
    twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// TwilioConfig berisi konfigurasi Twilio
type TwilioConfig struct {
    AccountSID  string
    AuthToken   string
    FromNumber  string
}

// TwilioProvider adalah implementasi provider SMS menggunakan Twilio
type TwilioProvider struct {
    client     *twilio.RestClient
    fromNumber string
}

// NewTwilioProvider membuat instance baru TwilioProvider
func NewTwilioProvider(config TwilioConfig) *TwilioProvider {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: config.AccountSID,
        Password: config.AuthToken,
    })

    return &TwilioProvider{
        client:     client,
        fromNumber: config.FromNumber,
    }
}

// SendSMS mengirim pesan SMS menggunakan Twilio
func (p *TwilioProvider) SendSMS(to, message string) error {
    params := &twilioApi.CreateMessageParams{
        To:   &to,
        From: &p.fromNumber,
        Body: &message,
    }

    _, err := p.client.Api.CreateMessage(params)
    return err
}
```

### Langkah 4: Registrasi Provider dalam Aplikasi Anda

```go
import (
    "github.com/kreasimaju/auth"
    "github.com/kreasimaju/auth/providers"
)

func main() {
    // Konfigurasi Twilio
    twilioConfig := providers.TwilioConfig{
        AccountSID:  "your-account-sid",
        AuthToken:   "your-auth-token",
        FromNumber:  "+1234567890",
    }

    // Buat provider Twilio
    twilioProvider := providers.NewTwilioProvider(twilioConfig)

    // Registrasi provider
    auth.SetSMSProvider(twilioProvider)

    // Lanjutkan dengan inisialisasi auth normal...
}
```

## Integrasi dengan Zenziva (Indonesia)

### Langkah 1: Buat Akun Zenziva
Daftar di [Zenziva](https://www.zenziva.id/) dan dapatkan:
- User Key
- API Key

### Langkah 2: Implementasi Provider

Buat file `providers/zenziva.go`:

```go
package providers

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
)

// ZenzivaConfig berisi konfigurasi Zenziva
type ZenzivaConfig struct {
    UserKey string
    APIKey  string
}

// ZenzivaProvider adalah implementasi provider SMS menggunakan Zenziva
type ZenzivaProvider struct {
    userKey string
    apiKey  string
    baseURL string
}

// NewZenzivaProvider membuat instance baru ZenzivaProvider
func NewZenzivaProvider(config ZenzivaConfig) *ZenzivaProvider {
    return &ZenzivaProvider{
        userKey: config.UserKey,
        apiKey:  config.APIKey,
        baseURL: "https://console.zenziva.net/masking/api/sendSMS/",
    }
}

// SendSMS mengirim pesan SMS menggunakan Zenziva
func (p *ZenzivaProvider) SendSMS(to, message string) error {
    // Buat URL endpoint
    endpoint := p.baseURL

    // Buat form data
    formData := url.Values{}
    formData.Set("userkey", p.userKey)
    formData.Set("passkey", p.apiKey)
    formData.Set("to", to)
    formData.Set("message", message)

    // Kirim request
    resp, err := http.PostForm(endpoint, formData)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Baca response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // Parse response JSON
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return err
    }

    // Periksa status
    if status, ok := result["status"].(string); ok && status != "success" {
        return fmt.Errorf("zenziva error: %s", result["text"])
    }

    return nil
}
```

### Langkah 4: Registrasi Provider dalam Aplikasi Anda

```go
import (
    "github.com/kreasimaju/auth"
    "github.com/kreasimaju/auth/providers"
)

func main() {
    // Konfigurasi Zenziva
    zenzivaConfig := providers.ZenzivaConfig{
        UserKey: "your-user-key",
        APIKey:  "your-api-key",
    }

    // Buat provider Zenziva
    zenzivaProvider := providers.NewZenzivaProvider(zenzivaConfig)

    // Registrasi provider
    auth.SetSMSProvider(zenzivaProvider)

    // Lanjutkan dengan inisialisasi auth normal...
}
```

## Menguji Integrasi SMS

Untuk menguji integrasi SMS, Anda dapat menggunakan kode berikut:

```go
func TestSMSIntegration() {
    // Kirim OTP melalui SMS
    otpCode, err := auth.RequestOTPLogin("081234567890", "sms", "ID")
    if err != nil {
        panic(err)
    }

    fmt.Printf("OTP code sent: %s\n", otpCode.Code)
}
```

## Troubleshooting

### Masalah Umum

1. **SMS Tidak Terkirim**
   - Periksa kredensi API Anda
   - Pastikan format nomor telepon benar (+62812345678 untuk Indonesia)
   - Cek saldo/kuota SMS Anda

2. **Error Format Nomor**
   - Gunakan fungsi `utils.FormatPhoneNumber` untuk memastikan format yang benar
   - Tentukan kode negara default yang sesuai (misal "ID" untuk Indonesia)

3. **Rate Limiting**
   - Beberapa provider memiliki batasan jumlah SMS yang dapat dikirim per menit/jam
   - Implementasikan mekanisme throttling jika diperlukan

### Implementasi Throttling

Anda dapat menambahkan throttling sederhana untuk mencegah pengiriman OTP berlebihan:

```go
// ThrottledSMSProvider membungkus provider SMS dengan throttling
type ThrottledSMSProvider struct {
    provider  SMSProvider
    lastSent  map[string]time.Time
    cooldown  time.Duration
    mu        sync.Mutex
}

func NewThrottledSMSProvider(provider SMSProvider, cooldownSeconds int) *ThrottledSMSProvider {
    return &ThrottledSMSProvider{
        provider:  provider,
        lastSent:  make(map[string]time.Time),
        cooldown:  time.Duration(cooldownSeconds) * time.Second,
        mu:        sync.Mutex{},
    }
}

func (p *ThrottledSMSProvider) SendSMS(to, message string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    // Cek apakah nomor dalam cooldown
    if lastTime, exists := p.lastSent[to]; exists {
        if time.Since(lastTime) < p.cooldown {
            return fmt.Errorf("Terlalu banyak permintaan. Coba lagi dalam %v detik", 
                p.cooldown.Seconds() - time.Since(lastTime).Seconds())
        }
    }

    // Kirim SMS
    err := p.provider.SendSMS(to, message)
    if err != nil {
        return err
    }

    // Catat waktu pengiriman
    p.lastSent[to] = time.Now()
    return nil
}
```

## Provider SMS Lokal Indonesia

Berikut adalah daftar beberapa provider SMS lokal Indonesia yang populer:

1. **Zenziva** - https://www.zenziva.id/
2. **Raja SMS** - https://rajasms.net/
3. **Infobip** - https://www.infobip.com/id/
4. **VoiceGate** - https://voicegate.id/
5. **Twilio** - https://www.twilio.com (Internasional tapi populer di Indonesia)

Pilihlah provider yang sesuai dengan kebutuhan dan budget Anda. Pertimbangkan faktor seperti:
- Harga per SMS
- Kecepatan pengiriman
- Jangkauan operator seluler
- Fitur tambahan (laporan pengiriman, API yang mudah, dll)
- Dukungan pelanggan 