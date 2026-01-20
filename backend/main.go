package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"medislink-backend/config"
	"medislink-backend/jobs"
	"medislink-backend/routes"
)

func main() {
	// In production (e.g., Railway), environment variables are injected by the platform
	// and a local .env file may not exist. Don't crash if it's missing.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; continuing with environment variables")
	}

	config.ConnectDB()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if v := os.Getenv("CORS_ALLOW_ORIGINS"); v != "" {
				return v
			}
			return "http://localhost:5173"
		}(),

		AllowHeaders: "Origin, Content-Type, Accept, Authorization",

		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	jobs.StartCronJob()

	routes.CategoryRoutes(app)
	routes.InventoryLogRoutes(app)
	routes.AddRoutes(app)
	routes.NewsRoutes(app)
	routes.NotificationRoutes(app)
	routes.ReviewRoutes(app)
	routes.LoanRoutes(app)
	routes.ReturnRoutes(app)
	routes.ToolRoutes(app)
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.DonationRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
