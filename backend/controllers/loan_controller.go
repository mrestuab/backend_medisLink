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
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateLoan(c *fiber.Ctx) error {
	var loan models.Loan

	if err := c.BodyParser(&loan); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	loan.UserID = userID

	if loan.Quantity <= 0 {
		loan.Quantity = 1
	}
	if loan.MedicalCondition == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Kondisi medis wajib diisi"})
	}

	loan.ID = primitive.NewObjectID()
	loan.LoanDate = time.Now().Format("2006-01-02")

	loan.Status = "pending"

	coll := config.DB.Collection("loans")
	_, err := coll.InsertOne(context.Background(), loan)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create loan"})
	}

	toolObjID, _ := primitive.ObjectIDFromHex(loan.ToolID)

	toolColl := config.DB.Collection("medical_tools")
	_, err = toolColl.UpdateOne(context.Background(),
		bson.M{"_id": toolObjID},
		bson.M{"$inc": bson.M{"stock": -loan.Quantity}},
	)

	if err != nil {
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Permintaan peminjaman berhasil dibuat",
		"loan_id": loan.ID,
	})
}

func GetAllLoans(c *fiber.Ctx) error {
	coll := config.DB.Collection("loans")
	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.D{
			{Key: "user_object_id", Value: bson.D{{Key: "$toObjectId", Value: "$user_id"}}},
			{Key: "tool_object_id", Value: bson.D{{Key: "$toObjectId", Value: "$tool_id"}}},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "user_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "user_detail"},
		}}},

		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user_detail"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "medical_tools"},
			{Key: "localField", Value: "tool_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tool_detail"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$tool_detail"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch loans"})
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode loans"})
	}

	return c.JSON(results)
}

func CompleteLoan(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	coll := config.DB.Collection("loans")
	_, err = coll.UpdateOne(context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": "selesai"}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update loan"})
	}

	return c.JSON(fiber.Map{"message": "Loan marked as completed"})
}

func GetMyLoans(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	coll := config.DB.Collection("loans")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},

		{{Key: "$addFields", Value: bson.D{
			{Key: "tool_object_id", Value: bson.D{{Key: "$toObjectId", Value: "$tool_id"}}},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "medical_tools"},
			{Key: "localField", Value: "tool_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tool_detail"},
		}}},

		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$tool_detail"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		{{Key: "$sort", Value: bson.D{{Key: "loan_date", Value: -1}}}},
	}

	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil riwayat pinjaman"})
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membaca data"})
	}

	if results == nil {
		return c.JSON([]bson.M{})
	}

	return c.JSON(results)
}

func UpdateLoanStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	type UpdateInput struct {
		Status string `json:"status"`
	}
	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	loanColl := config.DB.Collection("loans")
	toolColl := config.DB.Collection("medical_tools")

	var loan models.Loan
	err = loanColl.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&loan)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Loan not found"})
	}

	if input.Status == "rejected" || input.Status == "completed" {
		if loan.Status != "completed" && loan.Status != "rejected" {
			toolObjID, _ := primitive.ObjectIDFromHex(loan.ToolID)
			toolColl.UpdateOne(context.Background(),
				bson.M{"_id": toolObjID},
				bson.M{"$inc": bson.M{"stock": loan.Quantity}},
			)
		}
	}

	_, err = loanColl.UpdateOne(context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": input.Status}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update status"})
	}

	var title, message, notifType string

	switch input.Status {
	case "approved":
		title = "Permintaan Disetujui!"
		message = "Alat medis Anda sudah siap. Silakan datang ke klinik untuk pengambilan barang."
		notifType = "success"
	case "rejected":
		title = "Permintaan Ditolak"
		message = "Maaf, permintaan peminjaman Anda ditolak. Silakan cek ketersediaan atau hubungi admin."
		notifType = "error"
	case "active":
		title = "Peminjaman Dimulai"
		message = "Alat telah diserahkan kepada Anda. Harap dijaga dengan baik selama masa peminjaman."
		notifType = "info"
	case "completed":
		title = "Peminjaman Selesai"
		message = "Terima kasih telah mengembalikan alat medis tepat waktu. Semoga sehat selalu!"
		notifType = "success"
	}

	if title != "" {
		CreateNotificationService(loan.UserID, title, message, notifType)
	}

	return c.JSON(fiber.Map{"message": "Status berhasil diperbarui"})
}
