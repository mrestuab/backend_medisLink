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

func CreateInventoryLog(c *fiber.Ctx) error {
	var log models.InventoryLog
	if err := c.BodyParser(&log); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.ID = primitive.NewObjectID()
	log.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	coll := config.DB.Collection("inventory_logs")
	_, err := coll.InsertOne(context.Background(), log)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create log"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Inventory log created successfully"})
}

func GetAllInventoryLogs(c *fiber.Ctx) error {
	coll := config.DB.Collection("inventory_logs")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}
	defer cursor.Close(context.Background())

	var logs []models.InventoryLog
	cursor.All(context.Background(), &logs)

	return c.JSON(logs)
}
