package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id"`
	InvoiceID      string             `json:"invoice_id"`
	PaymentMethod  string             `json:"payment_method" validate:"eq=CARD|eq=CASH|eq""`
	PaymentStatus  string             `json:"payment_status" validate:"Required,eq=PENDNG|eq=PAID"`
	PaymentDueDate time.Time          `json:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
