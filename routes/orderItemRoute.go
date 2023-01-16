package routes

import (
	controller "canteenApi/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(r *gin.Engine) {
	r.GET("/orderItems", controller.GetOrderItems())
	r.GET("/orderItems/:orderItemId", controller.GetOrderItem())
	r.GET("/orderItem-order:orderId", controller.GetOrderItemsByOrder())
	r.POST("/orderItems", controller.CreateOrderItem())
	r.PATCH("/orderItems/:orderItemId", controller.UpdateOrderItem())
}
