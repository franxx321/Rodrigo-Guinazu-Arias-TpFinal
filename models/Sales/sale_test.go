package Sales

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"resty.dev/v3"
	"testing"
)

func TestCreateSale(t *testing.T) {
	// Configuración del servidor mock
	mockHandler := http.NewServeMux()

	// Caso de usuario existente
	mockHandler.HandleFunc("/users/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "id": "123",
            "name": "Juan Perez",
            "address": "Calle 123",
            "nickname": "juanp"
        }`))
	})

	// Caso de usuario no existente
	mockHandler.HandleFunc("/users/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"user not found"}`))
	})

	mockServer := httptest.NewServer(mockHandler)
	defer mockServer.Close()

	// Configurar el cliente resty para usar el servidor mock
	client := resty.New()
	client.SetBaseURL(mockServer.URL)
	fmt.Println("URL del mock server ", mockServer.URL)

	os.Setenv("USER_SERVICE_URL", mockServer.URL)

	t.Run("Crear venta exitosa", func(t *testing.T) {
		storage := NewSaleStorage()
		service := NewSaleService(storage)

		// Verificar que el endpoint responde correctamente usando resty
		resp, err := client.R().Get("/users/123")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		sale, err := service.Create("123", 100.50)

		assert.NoError(t, err)
		assert.NotNil(t, sale)
		assert.Equal(t, "123", sale.UserId)
		assert.Equal(t, float32(100.50), sale.Amount)
		assert.Contains(t, []string{Pending, Aproved, Rejected}, sale.Status)
		assert.NotEmpty(t, sale.Id)
		assert.Equal(t, 1, sale.Version)
	})

	t.Run("Crear venta con usuario inexistente", func(t *testing.T) {
		storage := NewSaleStorage()
		service := NewSaleService(storage)

		// Verificar que el endpoint responde correctamente usando resty
		resp, err := client.R().Get("/users/999")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode())

		sale, err := service.Create("999", 100.50)

		assert.Error(t, err)
		assert.Nil(t, sale)
	})
}
