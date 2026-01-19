package routes

import (
	"medislink-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/api/auth")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// Forgot Password Routes
	api.Post("/forgot-password", controllers.RequestPasswordReset)
	api.Post("/verify-otp", controllers.VerifyOTP)
	api.Post("/reset-password", controllers.ResetPassword)
}
