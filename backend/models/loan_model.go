package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Loan struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID           string             `json:"user_id" bson:"user_id"`
	ToolID           string             `json:"tool_id" bson:"tool_id"`
	Quantity         int                `json:"quantity" bson:"quantity"`
	MedicalCondition string             `json:"medical_condition" bson:"medical_condition"`
	LoanDate         string             `json:"loan_date" bson:"loan_date"`
	ReturnDue        string             `json:"return_due" bson:"return_due"`
	Status           string             `json:"status" bson:"status"`
	Notes            string             `json:"notes" bson:"notes"`
}
