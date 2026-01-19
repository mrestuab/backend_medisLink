package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InventoryLogRoutes(app *fiber.App) {
	api := app.Group("/api/inventory-logs")
	api.Use(middlewares.JWTProtected())

	api.Post("/", controllers.CreateInventoryLog)
	api.Get("/", controllers.GetAllInventoryLogs)
}
