package main

import (
	"Rodrigo-Guinazu-Arias-TpFinal/api"
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"Rodrigo-Guinazu-Arias-TpFinal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	//TOASK tenemos que hacer otro programa para interactuar con este programa
	//	como si este estuviera desplegado en otro lado o hacemos to do junto
	saleService := Sales.NewSaleService(Sales.NewSaleStorage())
	userService := users.NewUserService(users.NewUserStorage())
	utils.InitSystem(saleService, userService)
	r := gin.Default()
	api.InitRoutes(r, userService, saleService)
	if err := r.Run(":1234"); err != nil {
		panic(fmt.Errorf("error trying to start server: %v", err))
	}
}
