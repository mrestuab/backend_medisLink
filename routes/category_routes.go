package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(app *fiber.App) {
	api := app.Group("/api/categories")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateCategory)
	api.Get("/", controllers.GetAllCategories)
}
