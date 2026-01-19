package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"
	"medislink-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi email tidak boleh kosong
	if user.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email tidak boleh kosong"})
	}

	userColl := config.DB.Collection("users")

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	err := userColl.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// Email sudah ada
		return c.Status(409).JSON(fiber.Map{"error": "Email sudah terdaftar"})
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(hash)
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.Role = "user"

	// Insert user baru
	_, err = userColl.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userColl := config.DB.Collection("users")

	var user models.User
	err := userColl.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"token": signedToken})
}

// generateOTP menghasilkan kode OTP 6 digit
func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	// Format menjadi 6 digit dengan leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// RequestPasswordReset - User request untuk reset password
func RequestPasswordReset(c *fiber.Ctx) error {
	type ResetRequest struct {
		Email string `json:"email"`
	}
	var req ResetRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi email tidak kosong
	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email tidak boleh kosong"})
	}

	// Cek apakah email terdaftar
	userColl := config.DB.Collection("users")
	var user models.User
	err := userColl.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		// Untuk keamanan, jangan beritahu user bahwa email tidak terdaftar
		// Kirim response sukses agar attacker tidak tau email mana yang terdaftar
		return c.JSON(fiber.Map{"message": "Jika email terdaftar, kode OTP akan dikirim"})
	}

	// Generate OTP 6 digit
	otp, err := generateOTP()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal generate OTP"})
	}

	// Simpan OTP ke database dengan expiry 10 menit
	resetColl := config.DB.Collection("password_resets")

	// Hapus OTP lama yang belum digunakan untuk email ini
	resetColl.DeleteMany(context.Background(), bson.M{
		"email": req.Email,
		"used":  false,
	})

	passwordReset := models.PasswordReset{
		ID:        primitive.NewObjectID(),
		Email:     req.Email,
		OTP:       otp,
		ExpiresAt: primitive.NewDateTimeFromTime(time.Now().Add(10 * time.Minute)),
		Used:      false,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = resetColl.InsertOne(context.Background(), passwordReset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan OTP"})
	}

	// Kirim email OTP menggunakan Resend
	err = utils.SendOTPEmail(req.Email, otp)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengirim email OTP"})
	}

	return c.JSON(fiber.Map{
		"message": "Kode OTP telah dikirim ke email Anda",
	})
}

// VerifyOTP - Verifikasi kode OTP yang diinput user
func VerifyOTP(c *fiber.Ctx) error {
	type VerifyRequest struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	var req VerifyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi input
	if req.Email == "" || req.OTP == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email dan OTP tidak boleh kosong"})
	}

	// Cari OTP di database
	resetColl := config.DB.Collection("password_resets")
	var passwordReset models.PasswordReset

	err := resetColl.FindOne(context.Background(), bson.M{
		"email": req.Email,
		"otp":   req.OTP,
		"used":  false,
	}).Decode(&passwordReset)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Kode OTP tidak valid atau sudah digunakan"})
	}

	// Cek apakah OTP sudah expired
	if time.Now().After(passwordReset.ExpiresAt.Time()) {
		return c.Status(400).JSON(fiber.Map{"error": "Kode OTP sudah kadaluarsa"})
	}

	return c.JSON(fiber.Map{
		"message": "Kode OTP valid",
		"valid":   true,
	})
}

// ResetPassword - Update password baru setelah OTP terverifikasi
func ResetPassword(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	var req ResetPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi input
	if req.Email == "" || req.OTP == "" || req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email, OTP, dan password baru tidak boleh kosong"})
	}

	// Validasi panjang password minimal
	if len(req.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password minimal 6 karakter"})
	}

	// Cari OTP di database
	resetColl := config.DB.Collection("password_resets")
	var passwordReset models.PasswordReset

	err := resetColl.FindOne(context.Background(), bson.M{
		"email": req.Email,
		"otp":   req.OTP,
		"used":  false,
	}).Decode(&passwordReset)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Kode OTP tidak valid atau sudah digunakan"})
	}

	// Cek apakah OTP sudah expired
	if time.Now().After(passwordReset.ExpiresAt.Time()) {
		return c.Status(400).JSON(fiber.Map{"error": "Kode OTP sudah kadaluarsa"})
	}

	// Hash password baru
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}

	// Update password di database users
	userColl := config.DB.Collection("users")
	_, err = userColl.UpdateOne(
		context.Background(),
		bson.M{"email": req.Email},
		bson.M{"$set": bson.M{"password": string(hash)}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengupdate password"})
	}

	// Tandai OTP sebagai sudah digunakan
	resetColl.UpdateOne(
		context.Background(),
		bson.M{"_id": passwordReset.ID},
		bson.M{"$set": bson.M{"used": true}},
	)

	return c.JSON(fiber.Map{
		"message": "Password berhasil direset",
	})
}
