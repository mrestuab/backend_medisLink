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
	// Load .env file only in development (optional in production)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	config.ConnectDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://www.medislink.web.id/",

		AllowHeaders: "Origin, Content-Type, Accept, Authorization",

		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	jobs.StartCronJob()

	// Initialize all routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
