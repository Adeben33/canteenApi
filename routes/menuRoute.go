package routes

import (
	controller "canteenApi/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(r *gin.Engine) {
	r.GET("/menu", controller.GetMenus())
	r.GET("/menu/:menuId", controller.GetMenu())
	r.POST("/menu", controller.CreateMenu())
	r.PATCH("/menu/:menuId", controller.UpdateMenu())

}
