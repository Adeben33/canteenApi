package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id"`
	OrderDate time.Time          `json:"order_date" validate:"required"`
	CreatedAT time.Time          `json:"created_at"`
	UpdatedAT time.Time          `json:"updated_at"`
	OrderId   string             `json:"order_id"`
	TableId   *string            `json:"table_id" validation:"required"`
}
