package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type InventoryLog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ToolID    string             `json:"tool_id" bson:"tool_id"`
	Action    string             `json:"action" bson:"action"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Timestamp string             `json:"timestamp" bson:"timestamp"`
	Note      string             `json:"note" bson:"note"`
}
