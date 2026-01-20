package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Add struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Link        string             `json:"link" bson:"link"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
}
