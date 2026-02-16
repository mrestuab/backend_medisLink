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
)

// CreateAdd godoc
// @Summary Create advertisement
// @Tags Adds
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Param link formData string false "Link"
// @Param image formData file false "Image"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/adds/ [post]
func CreateAdd(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("image")
	var imageUrl string
	if err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal membuka file gambar"})
		}
		defer file.Close()
		imageUrl, err = utils.UploadToCloudinary(file, "adds")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal upload gambar ke Cloudinary", "details": err.Error()})
		}
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	link := c.FormValue("link")
	if title == "" || description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title dan description wajib diisi"})
	}

	add := models.Add{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Description: description,
		Link:        link,
		ImageURL:    imageUrl,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	coll := config.DB.Collection("adds")
	_, err = coll.InsertOne(context.Background(), add)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create add"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Add created successfully"})
}

// GetAllAdds godoc
// @Summary Get all advertisements
// @Tags Adds
// @Produce json
// @Success 200 {array} models.Add
// @Router /api/adds/ [get]
func GetAllAdds(c *fiber.Ctx) error {
	coll := config.DB.Collection("adds")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch adds"})
	}
	defer cursor.Close(context.Background())

	var adds []models.Add
	cursor.All(context.Background(), &adds)

	return c.JSON(adds)
}

// UpdateAdd godoc
// @Summary Update advertisement
// @Tags Adds
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Add ID"
// @Param add body models.Add true "Add Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/adds/{id} [put]
func UpdateAdd(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var add models.Add
	if err := c.BodyParser(&add); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	coll := config.DB.Collection("adds")
	_, err = coll.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": add})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update add"})
	}

	return c.JSON(fiber.Map{"message": "Add updated successfully"})
}

// DeleteAdd godoc
// @Summary Delete advertisement
// @Tags Adds
// @Security BearerAuth
// @Produce json
// @Param id path string true "Add ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/adds/{id} [delete]
func DeleteAdd(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	coll := config.DB.Collection("adds")
	_, err = coll.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete add"})
	}

	return c.JSON(fiber.Map{"message": "Add deleted"})
}
