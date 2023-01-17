package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	FirstName     *string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName      *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Avater        *string            `json:"avater"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_Token"`
	CreatedAt     time.Time          `json:"created_at"`
	updatedAt     time.Time          `json:"updated_at"`
	UserId        string             `json:"user_id"`
}
