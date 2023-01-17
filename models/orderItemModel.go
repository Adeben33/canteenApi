package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	Quantity    string             `json:"quantity" validate:"required,eq=s|eq=m|eql" `
	UnitPrice   float64            `json:"unitPrice" validate:"required"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAT   time.Time          `json:"updated_at"`
	FoodId      string             `json:"food_id" validate:"required"`
	OrderItemid string             `json:"order_itemid"`
	OrderID     string             `json:"order_id"validate:"required"`
}
