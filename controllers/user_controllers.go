package controllers

import (
	"context"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"
	"medislink-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return config.DB.Collection("users")
}

// CreateUser godoc
// @Summary Create user (admin/manual)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /api/users [post]
func CreateUser(c *fiber.Ctx) error {
	c.Context().SetUserValue("Content-Type", "application/json")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := getUserCollection().InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(user)
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Router /api/users/ [get]
func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse users"})
	}

	return c.JSON(users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var user models.User
	err = getUserCollection().FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags Users
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "User ID"
// @Param name formData string false "Name"
// @Param phone formData string false "Phone"
// @Param address formData string false "Address"
// @Param nik formData string false "NIK"
// @Param foto_profile formData file false "Profile Photo"
// @Param foto_ktp formData file false "KTP Photo"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users/{id} [put]
func UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	coll := config.DB.Collection("users")
	updateFields := bson.M{}

	if name := c.FormValue("name"); name != "" {
		updateFields["name"] = name
	}
	if phone := c.FormValue("phone"); phone != "" {
		updateFields["phone"] = phone
	}
	if address := c.FormValue("address"); address != "" {
		updateFields["address"] = address
	}
	if nik := c.FormValue("nik"); nik != "" {
		updateFields["nik"] = nik
	}

	fileProfile, err := c.FormFile("foto_profile")
	if err == nil {
		file, _ := fileProfile.Open()
		defer file.Close()
		url, errUpload := utils.UploadToCloudinary(file, "profiles")
		if errUpload == nil {
			updateFields["foto_profile"] = url
		}
	}

	fileKTP, err := c.FormFile("foto_ktp")
	if err == nil {
		file, _ := fileKTP.Open()
		defer file.Close()
		url, errUpload := utils.UploadToCloudinary(file, "ktp")
		if errUpload == nil {
			updateFields["foto_ktp"] = url
		}
	}

	if len(updateFields) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Tidak ada data yang diubah"})
	}

	_, err = coll.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update profile"})
	}

	return c.JSON(fiber.Map{"message": "Profil berhasil diperbarui"})
}

// DeleteUser godoc
// @Summary Delete user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	_, err = getUserCollection().DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted"})
}
