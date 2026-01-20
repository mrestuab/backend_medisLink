package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func NewsRoutes(app *fiber.App) {
	api := app.Group("/api/news")

	api.Get("/", controllers.GetAllNews)

	protected := api.Group("/")
	protected.Use(middlewares.JWTProtected())

	protected.Post("/", controllers.CreateNews)
	protected.Put("/:id", controllers.UpdateNews)
	protected.Delete("/:id", controllers.DeleteNews)
}
