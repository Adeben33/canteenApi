package routes

import "github.com/gin-gonic/gin"

func MenuRoutes(r *gin.Engine) {
	r.GET("/menu", controller.GetMenu())
	r.GET("/menu/:menuId", controller.GetMenu())
	r.POST("/menu", controller.CreateMenu())
	r.PATCH("/menu/:menuId", controller.UpdateMenu())

}
