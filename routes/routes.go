package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes all routes
func SetupRoutes(app *fiber.App) {
	AuthRoutes(app)
	UserRoutes(app)
	ToolRoutes(app)
	CategoryRoutes(app)
	DonationRoutes(app)
	LoanRoutes(app)
	NewsRoutes(app)
	NotificationRoutes(app)
	ReturnRoutes(app)
	ReviewRoutes(app)
	InventoryLogRoutes(app)
	AddRoutes(app)
}

// AuthRoutes handles authentication routes
func AuthRoutes(app *fiber.App) {
	api := app.Group("/api/auth")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// Forgot Password Routes
	api.Post("/forgot-password", controllers.RequestPasswordReset)
	api.Post("/verify-otp", controllers.VerifyOTP)
	api.Post("/reset-password", controllers.ResetPassword)
}

// UserRoutes handles user-related routes
func UserRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to MedisLink API 🚀")
	})

	api.Post("/users", controllers.CreateUser)

	protected := api.Group("/users", middlewares.JWTProtected())

	protected.Get("/", controllers.GetAllUsers)
	protected.Get("/:id", controllers.GetUserByID)
	protected.Put("/:id", controllers.UpdateProfile)
	protected.Delete("/:id", controllers.DeleteUser)
}

// ToolRoutes handles tool-related routes
func ToolRoutes(app *fiber.App) {
	api := app.Group("/api/tools")
	api.Get("/", controllers.GetAllTools)
	api.Get("/:id", controllers.GetToolByID)

	protected := api.Group("/")
	protected.Use(middlewares.JWTProtected())

	protected.Post("/", controllers.CreateTool)
	protected.Put("/:id", controllers.UpdateTool)
	protected.Delete("/:id", controllers.DeleteTool)
}

// CategoryRoutes handles category-related routes
func CategoryRoutes(app *fiber.App) {
	api := app.Group("/api/categories")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateCategory)
	api.Get("/", controllers.GetAllCategories)
}

// DonationRoutes handles donation-related routes
func DonationRoutes(app *fiber.App) {
	api := app.Group("/api/donations")

	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateDonation)
	api.Get("/history", controllers.GetUserDonations)
	api.Get("/", controllers.GetAllDonations)
	api.Put("/:id/approve", controllers.ApproveDonation)
	api.Put("/:id/receive", controllers.ReceiveDonation)
}

// LoanRoutes handles loan-related routes
func LoanRoutes(app *fiber.App) {
	api := app.Group("/api/loans")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateLoan)
	api.Get("/", controllers.GetAllLoans)
	api.Put("/:id/complete", controllers.CompleteLoan)
	api.Get("/my", controllers.GetMyLoans)
	api.Put("/:id/status", controllers.UpdateLoanStatus)
}

// NewsRoutes handles news-related routes
func NewsRoutes(app *fiber.App) {
	api := app.Group("/api/news")

	api.Get("/", controllers.GetAllNews)
	api.Get("/:id", controllers.GetNewsByID)

	protected := api.Group("/")
	protected.Use(middlewares.JWTProtected())

	protected.Post("/", controllers.CreateNews)
	protected.Put("/:id", controllers.UpdateNews)
	protected.Delete("/:id", controllers.DeleteNews)
}

// NotificationRoutes handles notification-related routes
func NotificationRoutes(app *fiber.App) {
	api := app.Group("/api/notifications")
	api.Use(middlewares.JWTProtected())

	api.Get("/", controllers.GetMyNotifications)
	api.Put("/:id/read", controllers.MarkAsRead)
}

// ReturnRoutes handles return-related routes
func ReturnRoutes(app *fiber.App) {
	api := app.Group("/api/returns")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateReturn)
	api.Get("/", controllers.GetAllReturns)
}

// ReviewRoutes handles review-related routes
func ReviewRoutes(app *fiber.App) {
	api := app.Group("/api/reviews")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateReview)
	api.Get("/:tool_id", controllers.GetToolReviews)
}

// InventoryLogRoutes handles inventory log-related routes
func InventoryLogRoutes(app *fiber.App) {
	api := app.Group("/api/inventory-logs")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateInventoryLog)
	api.Get("/", controllers.GetAllInventoryLogs)
}

// AddRoutes handles advertisement-related routes
func AddRoutes(app *fiber.App) {
	api := app.Group("/api/adds")

	api.Get("/", controllers.GetAllAdds)

	protected := api.Group("/")
	protected.Use(middlewares.JWTProtected())

	protected.Post("/", controllers.CreateAdd)
	protected.Put("/:id", controllers.UpdateAdd)
	protected.Delete("/:id", controllers.DeleteAdd)
}
