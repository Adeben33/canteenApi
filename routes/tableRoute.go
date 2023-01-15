package routes

import (
	controller "canteenApi/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(r *gin.Engine) {
	r.GET("/tables", controller.GetTables())
	r.GET("/tables/:tableId", controller.GetTables())
	r.POST("/table", controller.CreateTable())
	r.PATCH("/tables/:tableId", controller.UpdateTable())
}
