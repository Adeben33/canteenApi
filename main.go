package main

import (
	"canteenApi/database"
	"canteenApi/middleware"
	"canteenApi/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "Food")

//database.OpenCollection(database.Client, "Food")

func main() {
	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)
	router.Run(":" + port)

}
