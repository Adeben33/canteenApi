package routes

import (
	controller "canteenApi/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(r *gin.Engine) {
	r.GET("/foods", controller.GetFoods())
	r.GET("/foods/:foodId", controller.GetFood())
	r.POST("/foods", Controller.CreateFood())
	r.PATCH("/foods/:foodId", controller.UpdateFood())
}
