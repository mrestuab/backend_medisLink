package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	ToolID    string             `json:"tool_id" bson:"tool_id"`
	Rating    int                `json:"rating" bson:"rating"`
	Comment   string             `json:"comment" bson:"comment"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
