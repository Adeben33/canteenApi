package routes

import "github.com/gin-gonic/gin"

func OrderItemRoutes(r *gin.Engine) {
	r.GET("/orderItems", controller.GetOrdersItem())
	r.GET("/orderItems/:orderItemId", controller.GetOrderItem())
	r.POST("/orderItems", controller.CreateOrder())
	r.PATCH("/orderItems/:orderItemId", controller.UpdateOrder())
}
