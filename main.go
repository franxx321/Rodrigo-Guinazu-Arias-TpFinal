package main

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"Rodrigo-Guinazu-Arias-TpFinal/utils"
	"fmt"
	"resty.dev/v3"
)

func main() {

	//TOASK tenemos que hacer otro programa para interactuar con este programa
	//	como si este estuviera desplegado en otro lado o hacemos to do junto
	saleService := Sales.NewSaleService(Sales.NewSaleStorage())
	userService := users.NewUserService(users.NewLocalStorage())

	var parar = false
	var opcion = 0
	for !parar {
		for opcion < 1 || opcion > 5 {
			fmt.Println("Ingrese una opcion: " +
				"\n 1 Precargar Sistema \n " +
				"2 Crear Venta \n " +
				"3 Modificar Venta" +
				"\n 4 Mostrar Ventas por usuario \n" +
				" 5 Salir")
			_, err2 := fmt.Scanf("%d", &opcion)
			if err2 != nil {

			}
		}

		switch opcion {
		//TOIMPLEMENT poner todas las opciones que el sistema tiene que manejar
		//	si es que se hacen en este sistema y no en otro
		case 1:
			utils.InitSystem(saleService, userService)
			fmt.Println("Sistema precargado exitosamente!")
		case 2:
			fmt.Println("Ingrese el id del usuario: ")
			var id string
			fmt.Scanf("%s", &id)
			fmt.Println("Ingrese el monto de la venta: ")
			var monto float32
			fmt.Scanf("%f", &monto)
			var request struct {
				UserId string  `json:"userId"`
				Amount float32 `json:"amount"`
			}
			request.Amount = monto
			request.UserId = id
			req := resty.New()
			req.R().
				SetBody(request).Post("http://localhost:1234/sales")
		case 3:
			fmt.Println("Ingrese el id de la venta: ")
			var id string
			fmt.Scanf("%s", &id)
			fmt.Println("Ingrese el nuevo estado de la venta: ")
			var estado string
			fmt.Scanf("%s", &estado)
			var request struct {
				Status string `json:"status"`
			}
			req := resty.New()
			request.Status = estado
			req.R().SetQueryParam("id", id).
				SetBody(request).
				Patch("http://localhost:1234/sales")
		case 4:

		case 5:
			parar = true
		}
		if !parar {
			var auxParar int
			fmt.Println("Desea realizar otra operacion? 1 Si 2 No")
			fmt.Scanf("%d", &auxParar)
			if auxParar == 2 {
				parar = true
			} else {
				opcion = 0
			}
		}
	}
}
