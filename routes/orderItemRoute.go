package routes

import (
	controller "canteenApi/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(r *gin.Engine) {
	r.GET("/orderItems", controller.GetOrdersItem())
	r.GET("/orderItems/:orderItemId", controller.GetOrderItem())
	r.GET("/orderItem-order:orderId", controller.GetOrderItemsByOrder())
	r.POST("/orderItems", controller.CreateOrder())
	r.PATCH("/orderItems/:orderItemId", controller.UpdateOrder())
}
