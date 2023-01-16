package routes

import (
	controller "canteenApi/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(r *gin.Engine) {
	r.GET("/invoices", controller.GetInvoices())
	r.GET("/invoices/:invoiceId", controller.GetInvoice())
	r.POST("/invoices", controller.CreateInvoice())
	r.PATCH("/invoices/:invoiceId", controller.UpdateInvoice())

}
