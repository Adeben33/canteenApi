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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur listing the menu"})
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menuId := c.Param("menuId")
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"food_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching the food item"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		//Binding The request json into menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Binding error"})
			return
		}
		//	Validate the struct
		validateErr := validate.Struct(menu)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuID = menu.ID.Hex()
		result, InsertErr := menuCollection.InsertOne(ctx, menu)
		if InsertErr != nil {
			msg := fmt.Sprintf("Menu Item was not created")
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		menuId := c.Param("menuId")
		filter := bson.M{"menuId": menuId}

		var updateObj primitive.D
		if !(menu.StartDate.IsZero() && menu.EndDate.IsZero()) {
			if !inTimeSpan(menu.StartDate, menu.EndDate, time.Now()) {
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}
			updateObj = append(updateObj, bson.E{"start_Date", menu.StartDate})
			updateObj = append(updateObj, bson.E{"end_date", menu.EndDate})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{"name", menu.Name})
			}
			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{"category", menu.Category})
			}
			menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC822Z))
			updateObj = append(updateObj, bson.E{"update_at", menu.UpdatedAt})
			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			result, err := menuCollection.UpdateOne(
				ctx, filter, bson.D{
					{"$set", updateObj},
				},
				&opt)
			if err != nil {
				msg := "menu Update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}
			defer cancel()
			c.JSON(http.StatusOK, result)
		}
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
