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

// CreateReturn godoc
// @Summary Create return record
// @Tags Returns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param return body models.Return true "Return Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/returns/ [post]
func CreateReturn(c *fiber.Ctx) error {
	var ret models.Return
	if err := c.BodyParser(&ret); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	ret.ID = primitive.NewObjectID()
	ret.ReturnDate = time.Now().Format("2006-01-02")

	coll := config.DB.Collection("returns")
	_, err := coll.InsertOne(context.Background(), ret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create return"})
	}

	loanColl := config.DB.Collection("loans")
	_, _ = loanColl.UpdateOne(context.Background(),
		bson.M{"_id": ret.LoanID},
		bson.M{"$set": bson.M{"status": "selesai"}},
	)

	loan := models.Loan{}
	_ = loanColl.FindOne(context.Background(), bson.M{"_id": ret.LoanID}).Decode(&loan)
	toolColl := config.DB.Collection("medical_tools")
	_, _ = toolColl.UpdateOne(context.Background(),
		bson.M{"_id": loan.ToolID},
		bson.M{"$inc": bson.M{"stock": 1}},
	)

	return c.Status(201).JSON(fiber.Map{"message": "Return recorded successfully"})
}

// GetAllReturns godoc
// @Summary Get all returns
// @Tags Returns
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Return
// @Failure 401 {object} map[string]string
// @Router /api/returns/ [get]
func GetAllReturns(c *fiber.Ctx) error {
	coll := config.DB.Collection("returns")

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch returns"})
	}
	defer cursor.Close(context.Background())

	var returns []models.Return
	cursor.All(context.Background(), &returns)

	return c.JSON(returns)
}
