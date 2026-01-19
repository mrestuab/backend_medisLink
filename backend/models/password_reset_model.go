package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PasswordReset struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	OTP       string             `json:"otp" bson:"otp"`
	ExpiresAt primitive.DateTime `json:"expires_at" bson:"expires_at"`
	Used      bool               `json:"used" bson:"used"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}
