package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Return struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LoanID           string             `json:"loan_id" bson:"loan_id"`
	ReturnDate       string             `json:"return_date" bson:"return_date"`
	ConditionOnReturn string            `json:"condition_on_return" bson:"condition_on_return"`
	Penalty          int                `json:"penalty" bson:"penalty"`
}
