package controllers

import (
	"context"

	"medislink-backend/config"
	"medislink-backend/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(c *fiber.Ctx) error {
	var cat models.Category
	if err := c.BodyParser(&cat); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	cat.ID = primitive.NewObjectID()

	coll := config.DB.Collection("categories")
	_, err := coll.InsertOne(context.Background(), cat)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create category"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Category created successfully"})
}

func GetAllCategories(c *fiber.Ctx) error {
	coll := config.DB.Collection("categories")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}
	defer cursor.Close(context.Background())

	var categories []models.Category
	cursor.All(context.Background(), &categories)

	return c.JSON(categories)
}

