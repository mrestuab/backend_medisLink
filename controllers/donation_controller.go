package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"medislink-backend/config"
	"medislink-backend/models"
	"medislink-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ApproveRequest struct {
	Condition string `json:"condition"`
}

func CreateDonation(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Foto barang wajib diupload"})
	}
	file, _ := fileHeader.Open()
	defer file.Close()

	imageUrl, err := utils.UploadToCloudinary(file, "donations")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal upload gambar"})
	}

	qty, _ := strconv.Atoi(c.FormValue("quantity"))

	donation := models.Donation{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		ToolName:      c.FormValue("tool_name"),
		Category:      c.FormValue("category"),
		Description:   c.FormValue("description"),
		PickupAddress: c.FormValue("pickup_address"),
		PickupDate:    c.FormValue("pickup_date"),
		Quantity:      qty,
		ImageURL:      imageUrl,
		Status:        "pending",
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	coll := config.DB.Collection("donations")
	_, err = coll.InsertOne(context.Background(), donation)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengajukan donasi"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Terima kasih! Donasi Anda sedang direview."})
}

func GetAllDonations(c *fiber.Ctx) error {
	coll := config.DB.Collection("donations")

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user_info",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$user_info",
			"preserveNullAndEmptyArrays": true,
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"tool_name":      1,
			"category":       1,
			"description":    1,
			"image_url":      1,
			"quantity":       1,
			"pickup_address": 1,
			"pickup_date":    1,
			"status":         1,
			"created_at":     1,
			"user_detail": bson.M{
				"name":  "$user_info.name",
				"phone": "$user_info.phone",
				"email": "$user_info.email",
			},
		}}},
		{{Key: "$sort", Value: bson.M{"created_at": -1}}},
	}

	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal load data"})
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal parsing data"})
	}

	return c.JSON(results)
}

func GetUserDonations(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	coll := config.DB.Collection("donations")

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := coll.Find(context.Background(), bson.M{"user_id": userID}, opts)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	defer cursor.Close(context.Background())

	var donations []models.Donation
	if err = cursor.All(context.Background(), &donations); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal parsing data"})
	}

	if donations == nil {
		donations = []models.Donation{}
	}

	return c.JSON(donations)
}

func ApproveDonation(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	var req ApproveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format data salah"})
	}
	if req.Condition == "" {
		req.Condition = "Layak Pakai"
	}

	collDonation := config.DB.Collection("donations")
	var donation models.Donation
	err := collDonation.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&donation)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data donasi tidak ditemukan"})
	}

	if donation.Status == "approved" {
		return c.Status(400).JSON(fiber.Map{"error": "Donasi ini sudah selesai diproses (Stok sudah masuk)"})
	}

	collTools := config.DB.Collection("medical_tools")
	var existingTool models.MedicalTool

	filter := bson.M{
		"name":      donation.ToolName,
		"condition": req.Condition,
	}

	err = collTools.FindOne(context.Background(), filter).Decode(&existingTool)

	if err == nil {
		update := bson.M{
			"$inc": bson.M{"stock": donation.Quantity},
			"$set": bson.M{"updated_at": time.Now()},
		}
		collTools.UpdateOne(context.Background(), filter, update)
	} else {
		newTool := models.MedicalTool{
			ID:          primitive.NewObjectID(),
			Name:        donation.ToolName,
			CategoryID:  donation.Category,
			Description: fmt.Sprintf("%s (Sumber: Donasi)", donation.Description),
			ImageURL:    donation.ImageURL,
			Stock:       donation.Quantity,
			Status:      "tersedia",
			Condition:   req.Condition,
			Type:        "Donasi",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		collTools.InsertOne(context.Background(), newTool)
	}

	_, err = collDonation.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"status":         "approved",
			"admin_qc_notes": req.Condition,
		}},
	)

	return c.JSON(fiber.Map{
		"message": "Sukses! Barang validasi QC selesai dan stok bertambah.",
	})
}

func ReceiveDonation(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	coll := config.DB.Collection("donations")

	_, err := coll.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": "received"}},
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update status"})
	}

	return c.JSON(fiber.Map{"message": "Barang berhasil diterima oleh Admin. Lanjutkan ke QC."})
}
