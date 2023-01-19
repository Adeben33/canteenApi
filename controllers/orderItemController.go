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

type OrderItemPack struct {
	TableId    *string
	OrderItems []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurs while listing all item"})
			return
		}
		var allOrderItem []bson.M
		if err = result.All(ctx, allOrderItem); err != nil {
			log.Panic(err)
			return
		}
		c.JSON(http.StatusOK, allOrderItem)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("orderId")
		allOrderItem, err := ItemByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order item by orderid"})
			return
		}
		c.JSON(http.StatusOK, allOrderItem)
	}
}

func ItemByOrder(id string) (orderItem []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookUpStage := bson.D{{"$lookup", bson.D{{"$from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderstage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"perserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreginField", "table_id"}, {"as", "table"}}}}
	unWindTableStage := bson.D{{"unwind", bson.D{{"path", "table"}, {"perserveNullAndEmptyArrays", "true"}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"foodname", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"order_od", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		},
		},
	}
	groupStage := bson.D{{"$group", bson.D{{"order_id", "$order_id"}, {"table_id", "$table.id"}, {"table_number", "$table.table_number"}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", "1"}}}, {"order_items", bson.D{}}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "_id.table_number"},
			{"order_items", 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookUpStage,
		unwindStage,
		lookupOrderstage,
		unwindOrderStage,
		lookupTableStage,
		unWindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}
	if err = result.All(ctx, &orderItem); err != nil {
		panic(err)
	}
	defer cancel()
	return orderItem, err
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderItemId := c.Param("orderItemId")
		var orderItem models.Order

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItemPack OrderItemPack
		var order models.Order
		if err := c.BindJSON(orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC822Z))
		order.CreatedAT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC822Z))
		orderItemToBeInserted := []interface{}{}
		order.TableId = orderItemPack.TableId
		orderId := OrderItemOrderCreator(order)
		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderID = orderId
			validateErr := validate.Struct(orderItem)

			if validateErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemid = orderItem.ID.Hex()
			var num = toFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			orderItemToBeInserted = append(orderItemToBeInserted, orderItem)

		}
		insertedOrderItem, err := orderItemCollection.InsertMany(ctx, orderItemToBeInserted)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, insertedOrderItem)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItem models.OrderItem
		orderItemId := c.Param("orderItemId")
		filter := bson.M{"order_item_id": orderItemId}

		var updateObj primitive.D
		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, bson.E{"quantity", *&orderItem.UnitPrice})
		}
		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", *orderItem.Quantity})
		}
		if orderItem.FoodId != nil {
			updateObj = append(updateObj, bson.E{"food_id", *orderItem.FoodId})
		}
		orderItem.UpdatedAT, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.UpdatedAT})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateObj},
		},
			&opt)
		if err != nil {
			msg := fmt.Sprintf("Order item Update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		c.JSON(http.StatusOK, result)
	}
}
