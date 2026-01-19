package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func LoanRoutes(app *fiber.App) {
	api := app.Group("/api/loans")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateLoan)
	api.Get("/", controllers.GetAllLoans)
	api.Put("/:id/complete", controllers.CompleteLoan)
	api.Get("/my", controllers.GetMyLoans)
	api.Put("/:id/status", controllers.UpdateLoanStatus)
}
