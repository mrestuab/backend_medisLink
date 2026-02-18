package controllers

import (
	"context"
	"strconv"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"
	"medislink-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTool(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gambar wajib diupload"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file"})
	}
	defer file.Close()

	imageUrl, err := utils.UploadToCloudinary(file, "alat_medis")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal upload ke Cloudinary", "details": err.Error()})
	}

	stockStr := c.FormValue("stock")
	stockInt, _ := strconv.Atoi(stockStr)

	tool := models.MedicalTool{
		ID:          primitive.NewObjectID(),
		Name:        c.FormValue("name"),
		CategoryID:  c.FormValue("category_id"),
		Type:        c.FormValue("type"),
		Size:        c.FormValue("size"),
		Dimensions:  c.FormValue("dimensions"),
		WeightCap:   c.FormValue("weight_cap"),
		Description: c.FormValue("description"),
		Condition:   c.FormValue("condition"),
		Stock:       stockInt,
		Status:      "tersedia",
		ImageURL:    imageUrl,
	}

	coll := config.DB.Collection("medical_tools")
	_, err = coll.InsertOne(context.Background(), tool)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add tool to database"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Tool added successfully",
		"data":    tool,
	})
}

func GetAllTools(c *fiber.Ctx) error {
	coll := config.DB.Collection("medical_tools")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tools"})
	}
	defer cursor.Close(ctx)

	var tools []models.MedicalTool
	if err := cursor.All(ctx, &tools); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse data"})
	}

	return c.JSON(tools)
}

func GetToolByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	coll := config.DB.Collection("medical_tools")
	var tool models.MedicalTool
	err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&tool)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tool not found"})
	}

	return c.JSON(tool)
}

func UpdateTool(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Parse request body
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Add fields to update if they exist in request
	if name, ok := updateData["name"].(string); ok && name != "" {
		update["$set"].(bson.M)["name"] = name
	}
	if categoryID, ok := updateData["category_id"].(string); ok {
		update["$set"].(bson.M)["category_id"] = categoryID
	}
	if toolType, ok := updateData["type"].(string); ok {
		update["$set"].(bson.M)["type"] = toolType
	}
	if size, ok := updateData["size"].(string); ok {
		update["$set"].(bson.M)["size"] = size
	}
	if dimensions, ok := updateData["dimensions"].(string); ok {
		update["$set"].(bson.M)["dimensions"] = dimensions
	}
	if weightCap, ok := updateData["weight_cap"].(string); ok {
		update["$set"].(bson.M)["weight_cap"] = weightCap
	}
	if description, ok := updateData["description"].(string); ok {
		update["$set"].(bson.M)["description"] = description
	}
	if condition, ok := updateData["condition"].(string); ok {
		update["$set"].(bson.M)["condition"] = condition
	}
	if status, ok := updateData["status"].(string); ok {
		update["$set"].(bson.M)["status"] = status
	}

	// Handle stock (can be float64 from JSON)
	if stock, ok := updateData["stock"]; ok {
		switch v := stock.(type) {
		case float64:
			update["$set"].(bson.M)["stock"] = int(v)
		case int:
			update["$set"].(bson.M)["stock"] = v
		case string:
			if stockInt, err := strconv.Atoi(v); err == nil {
				update["$set"].(bson.M)["stock"] = stockInt
			}
		}
	}

	coll := config.DB.Collection("medical_tools")
	result, err := coll.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update tool"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Tool not found"})
	}

	return c.JSON(fiber.Map{
		"message": "Tool updated successfully",
		"id":      id,
	})
}

func DeleteTool(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	coll := config.DB.Collection("medical_tools")
	_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete tool"})
	}

	return c.JSON(fiber.Map{"message": "Tool deleted"})
}
