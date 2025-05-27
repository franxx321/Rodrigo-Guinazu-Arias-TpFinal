// client_test.go

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// UserResponse representa la estructura de un usuario que esperamos recibir de la API.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	NickName  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

func main2() {
	client := resty.New()
	client.SetBaseURL("http://localhost:1234")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nSeleccione una operación:")
		fmt.Println("1. Crear usuario")
		fmt.Println("2. Obtener usuario")
		fmt.Println("3. Eliminar usuario")
		fmt.Println("4. Salir")
		fmt.Print("Opción: ")

		scanner.Scan()
		option := scanner.Text()

		switch option {
		case "1":
			createUser(client, scanner)
		case "2":
			getUser(client, scanner)
		case "3":
			deleteUser(client, scanner)
		case "4":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func createUser(client *resty.Client, scanner *bufio.Scanner) {
	fmt.Print("Nombre: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Dirección: ")
	scanner.Scan()
	address := scanner.Text()

	fmt.Print("Nickname: ")
	scanner.Scan()
	nickname := scanner.Text()

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"name":     name,
			"address":  address,
			"nickname": nickname,
		}).
		Post("/users")

	if err != nil {
		fmt.Printf("Error de conexión: %v\n", err)
		return
	}

	// Imprimir la respuesta completa
	fmt.Printf("\nCódigo de estado: %d\n", response.StatusCode())
	fmt.Printf("Cuerpo de la respuesta: %s\n", response.String())

	switch response.StatusCode() {
	case http.StatusCreated:
		fmt.Println("Usuario creado exitosamente")
	case http.StatusBadRequest:
		fmt.Println("Error en la solicitud: datos inválidos")
	case http.StatusInternalServerError:
		fmt.Println("Error interno del servidor")
	default:
		fmt.Printf("Código de respuesta inesperado: %d\n", response.StatusCode())
	}
}

func getUser(client *resty.Client, scanner *bufio.Scanner) {
	fmt.Print("ID del usuario: ")
	scanner.Scan()
	id := scanner.Text()

	response, err := client.R().
		Get("/users/" + id)

	if err != nil {
		fmt.Printf("Error de conexión: %v\n", err)
		return
	}

	// Imprimir la respuesta completa
	fmt.Printf("\nCódigo de estado: %d\n", response.StatusCode())
	fmt.Printf("Cuerpo de la respuesta: %s\n", response.String())

	switch response.StatusCode() {
	case http.StatusOK:
		fmt.Println("Usuario encontrado")
	case http.StatusNotFound:
		fmt.Println("Usuario no encontrado")
	case http.StatusInternalServerError:
		fmt.Println("Error interno del servidor")
	default:
		fmt.Printf("Código de respuesta inesperado: %d\n", response.StatusCode())
	}
}

func deleteUser(client *resty.Client, scanner *bufio.Scanner) {
	fmt.Print("ID del usuario: ")
	scanner.Scan()
	id := scanner.Text()

	response, err := client.R().
		Delete("/users/" + id)

	if err != nil {
		fmt.Printf("Error de conexión: %v\n", err)
		return
	}

	// Imprimir la respuesta completa
	fmt.Printf("\nCódigo de estado: %d\n", response.StatusCode())
	fmt.Printf("Cuerpo de la respuesta: %s\n", response.String())

	switch response.StatusCode() {
	case http.StatusNoContent:
		fmt.Println("Usuario eliminado exitosamente")
	case http.StatusNotFound:
		fmt.Println("Usuario no encontrado")
	case http.StatusInternalServerError:
		fmt.Println("Error interno del servidor")
	default:
		fmt.Printf("Código de respuesta inesperado: %d\n", response.StatusCode())
	}
}
