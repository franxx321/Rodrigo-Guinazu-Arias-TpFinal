package api

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TOIMPLEMENT toda la logica para la inicializacion de las urls
func InitRoutes(e *gin.Engine, userService *users.UserService, saleService *Sales.SaleService) {
	userHandler := UserHandler{userService}
	saleHandler := SaleHandler{saleService}

	e.POST("/users", userHandler.HandleUserCreate)
	e.GET("/users/:id", userHandler.HandleUserRead)
	//FIXME descomentar esa linea cuando joako arregle el error
	//e.PATCH("/users/:id", UserHandler.HandleUserUpdate)
	e.DELETE("/users/:id", userHandler.HandleUserDelete)
	e.POST("/sales", saleHandler.HandleSaleCreate)
	e.PATCH("/sales/:id", saleHandler.HandleSaleUpdate)
	e.GET("/sales", saleHandler.HandleSaleGetByUserStatus)
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
