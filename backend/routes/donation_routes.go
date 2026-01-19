package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func DonationRoutes(app *fiber.App) {
	api := app.Group("/api/donations")

	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateDonation)
	api.Get("/history", controllers.GetUserDonations)
	api.Get("/", controllers.GetAllDonations)
	api.Put("/:id/approve", controllers.ApproveDonation)
	api.Put("/:id/receive", controllers.ReceiveDonation)
}
