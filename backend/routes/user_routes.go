package routes

import (
	"medislink-backend/controllers"
	"medislink-backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

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
