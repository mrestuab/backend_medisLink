package controllers

import (
	"context"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"
	"medislink-backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateNews godoc
// @Summary Create news
// @Tags News
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Param image formData file false "Image"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/news/ [post]
func CreateNews(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	var imageUrl string
	if err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file gambar"})
		}
		defer file.Close()
		imageUrl, err = utils.UploadToCloudinary(file, "news")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal upload gambar ke Cloudinary", "details": err.Error()})
		}
	}

	title := c.FormValue("title")
	content := c.FormValue("content")
	if title == "" || content == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title dan content wajib diisi"})
	}

	news := models.News{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		ImageURL:  imageUrl,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	coll := config.DB.Collection("news")
	_, err = coll.InsertOne(context.Background(), news)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create news"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "News created successfully"})
}

// GetAllNews godoc
// @Summary Get all news
// @Tags News
// @Produce json
// @Success 200 {array} models.News
// @Router /api/news/ [get]
func GetAllNews(c *fiber.Ctx) error {
	coll := config.DB.Collection("news")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch news"})
	}
	defer cursor.Close(context.Background())

	var newsList []models.News
	cursor.All(context.Background(), &newsList)

	return c.JSON(newsList)
}

// GetNewsByID godoc
// @Summary Get news by ID
// @Tags News
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} models.News
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/news/{id} [get]
func GetNewsByID(c *fiber.Ctx) error {
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	coll := config.DB.Collection("news")

	var news models.News
	err = coll.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&news)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "News not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.JSON(news)
}

// UpdateNews godoc
// @Summary Update news
// @Tags News
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Param news body models.News true "News Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/news/{id} [put]
func UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var news models.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	coll := config.DB.Collection("news")
	_, err = coll.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": news})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update news"})
	}

	return c.JSON(fiber.Map{"message": "News updated successfully"})
}

// DeleteNews godoc
// @Summary Delete news
// @Tags News
// @Security BearerAuth
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/news/{id} [delete]
func DeleteNews(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	coll := config.DB.Collection("news")
	_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete news"})
	}

	return c.JSON(fiber.Map{"message": "News deleted"})
}
