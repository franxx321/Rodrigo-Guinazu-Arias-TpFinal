package main

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"Rodrigo-Guinazu-Arias-TpFinal/utils"
	"fmt"
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

		case 3:

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
