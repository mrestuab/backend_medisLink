package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AddRoutes(app *fiber.App) {
	api := app.Group("/api/adds")

	api.Get("/", controllers.GetAllAdds)

	protected := api.Group("/")
	protected.Use(middlewares.JWTProtected())

	protected.Post("/", controllers.CreateAdd)
	protected.Put("/:id", controllers.UpdateAdd)
	protected.Delete("/:id", controllers.DeleteAdd)
}
