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

// CreateReview godoc
// @Summary Create review
// @Tags Reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param review body models.Review true "Review Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/reviews/ [post]
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

// GetToolReviews godoc
// @Summary Get reviews by tool ID
// @Tags Reviews
// @Security BearerAuth
// @Produce json
// @Param tool_id path string true "Tool ID"
// @Success 200 {array} models.Review
// @Failure 401 {object} map[string]string
// @Router /api/reviews/{tool_id} [get]
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
