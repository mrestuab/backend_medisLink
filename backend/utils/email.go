package utils

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

// SendOTPEmail mengirim email berisi kode OTP untuk reset password
func SendOTPEmail(toEmail, otp string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY tidak ditemukan di environment variable")
	}

	client := resend.NewClient(apiKey)

	// HTML template untuk email
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 30px auto;
            background-color: #ffffff;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #14b8a6 0%%, #0891b2 100%%);
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            color: white;
            margin: 0;
            font-size: 28px;
        }
        .content {
            padding: 40px 30px;
        }
        .otp-box {
            background-color: #f0fdfa;
            border: 2px dashed #14b8a6;
            border-radius: 8px;
            padding: 20px;
            text-align: center;
            margin: 30px 0;
        }
        .otp-code {
            font-size: 36px;
            font-weight: bold;
            color: #0891b2;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
        }
        .warning {
            background-color: #fef2f2;
            border-left: 4px solid #ef4444;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        .warning p {
            margin: 0;
            color: #991b1b;
            font-size: 14px;
        }
        .footer {
            background-color: #f9fafb;
            padding: 20px;
            text-align: center;
            color: #6b7280;
            font-size: 12px;
        }
        .button {
            display: inline-block;
            padding: 12px 30px;
            background-color: #14b8a6;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 MedisLink</h1>
        </div>
        <div class="content">
            <h2 style="color: #1f2937; margin-top: 0;">Reset Password Anda</h2>
            <p style="color: #4b5563; line-height: 1.6;">
                Kami menerima permintaan untuk mereset password akun Anda. 
                Gunakan kode OTP berikut untuk melanjutkan proses reset password:
            </p>
            
            <div class="otp-box">
                <p style="margin: 0 0 10px 0; color: #6b7280; font-size: 14px;">Kode OTP Anda:</p>
                <div class="otp-code">%s</div>
                <p style="margin: 10px 0 0 0; color: #6b7280; font-size: 12px;">Berlaku selama 10 menit</p>
            </div>

            <div class="warning">
                <p>
                    <strong>⚠️ Penting:</strong> Jangan bagikan kode ini kepada siapapun. 
                    Tim MedisLink tidak akan pernah meminta kode OTP Anda.
                </p>
            </div>

            <p style="color: #4b5563; line-height: 1.6;">
                Jika Anda tidak meminta reset password, abaikan email ini. 
                Akun Anda tetap aman.
            </p>
        </div>
        <div class="footer">
            <p>© 2026 MedisLink - Sistem Manajemen Alat Kesehatan</p>
            <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
        </div>
    </div>
</body>
</html>
	`, otp)

	params := &resend.SendEmailRequest{
		From:    "MedisLink <onboarding@resend.dev>", // Ganti dengan domain terverifikasi Anda
		To:      []string{toEmail},
		Subject: "Reset Password - Kode OTP MedisLink",
		Html:    htmlContent,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %v", err)
	}

	fmt.Printf("Email OTP terkirim dengan ID: %s\n", sent.Id)
	return nil
}
