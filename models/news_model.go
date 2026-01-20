package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Content     string             `json:"content" bson:"content"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Author      string             `json:"author" bson:"author"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
}
