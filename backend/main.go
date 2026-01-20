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
	// Load .env file only in development (Railway akan pakai environment variables)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config.ConnectDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",

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

	log.Fatal(app.Listen(":" + port))
}
