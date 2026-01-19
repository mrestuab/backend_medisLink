package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

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
