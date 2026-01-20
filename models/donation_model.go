package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Donation struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	ToolName    string `json:"tool_name" bson:"tool_name"`
	Category    string `json:"category" bson:"category"`
	Description string `json:"description" bson:"description"`
	ImageURL    string `json:"image_url" bson:"image_url"`
	Quantity    int    `json:"quantity" bson:"quantity"`

	PickupAddress string `json:"pickup_address" bson:"pickup_address"`
	PickupDate    string `bson:"pickup_date" json:"pickup_date"`

	Status    string `json:"status" bson:"status"`
	CreatedAt string `json:"created_at" bson:"created_at"`

	UserDetail *UserResponse `json:"user_detail,omitempty" bson:"user_detail,omitempty"`
}
type UserResponse struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type ApproveRequest struct {
	Condition string `json:"condition"`
}
