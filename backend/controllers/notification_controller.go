package controllers

import (
	"context"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMyNotifications(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	coll := config.DB.Collection("notifications")

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := coll.Find(context.Background(), bson.M{"user_id": userID}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}
	defer cursor.Close(context.Background())

	var notifs []models.Notification
	if err = cursor.All(context.Background(), &notifs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode notifications"})
	}

	if notifs == nil {
		return c.JSON([]models.Notification{})
	}

	return c.JSON(notifs)
}

func MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	coll := config.DB.Collection("notifications")
	_, err = coll.UpdateOne(context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"is_read": true}},
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update"})
	}

	return c.JSON(fiber.Map{"message": "Notification marked as read"})
}

func CreateNotificationService(userID string, title string, message string, notifType string) {
	coll := config.DB.Collection("notifications")

	notif := models.Notification{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Type:      notifType,
	}

	go func() {
		coll.InsertOne(context.Background(), notif)
	}()
}
