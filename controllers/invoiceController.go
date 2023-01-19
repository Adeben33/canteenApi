package controllers

import (
	"canteenApi/database"
	"canteenApi/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type invoiceViewFormat struct {
	InvoiceId      string
	PaymentMethod  string
	OrderId        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			msg := fmt.Sprintf("Unable to fetch invoice")
			c.JSON(http.StatusInternalServerError, msg)
		}
		var allInvoices []bson.D
		resultErr := result.All(ctx, allInvoices)
		if resultErr != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		invoiceId := c.Param("invoiceId")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur while listing invoice item"})
		}
		var invoiceView invoiceViewFormat

		allOrderItem, err := ItemByOrder(invoice.InvoiceID)
		invoiceView.OrderId = invoice.InvoiceID
		invoiceView.PaymentDueDate = invoice.PaymentDueDate

		invoiceView.PaymentMethod = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}
		invoiceView.InvoiceId = invoice.InvoiceID
		invoiceView.PaymentStatus = *&invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItem[0]["paymentDue"]
		invoiceView.TableNumber = allOrderItem[0]["tableNumber"]
		invoiceView.OrderDetails = allOrderItem[0]["orderItem"]
		c.JSON(http.StatusOK, invoiceView)
	}

}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		//Binding The request json into invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Binding error"})
			return
		}
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderId}).Decode(&order)
		if err != nil {
			msg := fmt.Sprintf("message:order not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC822Z))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC822Z))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()
		validateErr := validate.Struct(invoice)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		result, insertErr := invoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil {
			msg := fmt.Sprintf("invoice item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		//Binding The request json into invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Binding error"})
			return
		}
		filter := bson.M{"invoice_id": invoice.InvoiceID}

		var updatedObj primitive.D
		if invoice.PaymentMethod != nil {
			updatedObj = append(updatedObj, bson.E{"payment_method", invoice.PaymentMethod})
		}
		if invoice.PaymentStatus != nil {
			updatedObj = append(updatedObj, bson.E{"payment_status", invoice.PaymentStatus})

		}
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC822Z))
		updatedObj = append(updatedObj, bson.E{"update_at", invoice.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		result, err := invoiceCollection.UpdateOne(
			ctx, filter, bson.D{
				{"$set", updatedObj},
			},
			&opt)
		if err != nil {
			msg := fmt.Sprintf("Invoice update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
