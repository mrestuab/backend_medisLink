package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Role        string             `json:"role" bson:"role"`
	Phone       string             `json:"phone" bson:"phone"`
	Address     string             `json:"address" bson:"address"`
	NIK         string             `json:"nik" bson:"nik"`
	FotoKTP     string             `json:"foto_ktp" bson:"foto_ktp"`
	FotoProfile string             `json:"foto_profile" bson:"foto_profile"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
}
