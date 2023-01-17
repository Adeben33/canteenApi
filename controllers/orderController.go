package controllers

import (
	"canteenApi/database"
	"context"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing food items"}
			)
		}

		var allOrders []bson.M

		if err := result.All(ctx, &allOrders); errr != nil {
			log.Panic(err)
		}
		c.JSON(http.StatusOK,allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	orderId := c.Param("orderId")
	var order models.O
	result,err = foodCollection.FindOne(ctx, bson.M{"order_id":orderId}, opt)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
