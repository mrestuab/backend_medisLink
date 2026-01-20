package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ReturnRoutes(app *fiber.App) {
	api := app.Group("/api/returns")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateReturn)
	api.Get("/", controllers.GetAllReturns)
}
