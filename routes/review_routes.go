package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ReviewRoutes(app *fiber.App) {
	api := app.Group("/api/reviews")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateReview)
	api.Get("/:tool_id", controllers.GetToolReviews)
}
