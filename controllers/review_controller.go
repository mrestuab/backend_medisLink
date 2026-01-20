package controllers

import (
	"context"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateReview(c *fiber.Ctx) error {
	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	coll := config.DB.Collection("reviews")
	_, err := coll.InsertOne(context.Background(), review)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save review"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Review added"})
}

func GetToolReviews(c *fiber.Ctx) error {
	toolID := c.Params("tool_id")
	coll := config.DB.Collection("reviews")

	cursor, err := coll.Find(context.Background(), bson.M{"tool_id": toolID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reviews"})
	}
	defer cursor.Close(context.Background())

	var reviews []models.Review
	cursor.All(context.Background(), &reviews)

	return c.JSON(reviews)
}
