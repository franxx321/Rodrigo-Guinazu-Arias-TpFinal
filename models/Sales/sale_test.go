package Sales

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
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

	t.Run("Crear venta exitosa usando resty", func(t *testing.T) {
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

	t.Run("Crear venta con usuario inexistente usando resty", func(t *testing.T) {
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

func TestUpdateSale(t *testing.T) {
	storage := NewSaleStorage()
	service := NewSaleService(storage)

	t.Run("Actualizar venta pendiente a aprobada", func(t *testing.T) {
		// Crear una venta en estado pendiente
		sale := &Sale{
			Id:        "sale-1",
			UserId:    "user-1",
			Amount:    100.0,
			Status:    Pending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		}
		storage.PutSale(sale)

		err := service.Update("sale-1", Aproved)
		assert.NoError(t, err)

		updatedSale, _ := storage.GetSale("sale-1")
		assert.Equal(t, Aproved, updatedSale.Status)
		assert.Equal(t, 2, updatedSale.Version)
	})

	t.Run("Actualizar venta con estado inválido", func(t *testing.T) {
		err := service.Update("sale-1", "estado_invalido")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidStatus, err)
	})

	t.Run("Actualizar venta inexistente", func(t *testing.T) {
		err := service.Update("venta-no-existe", Aproved)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("Actualizar a estado pendiente", func(t *testing.T) {
		err := service.Update("sale-1", Pending)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidTransition, err)
	})
}

func TestGetByUserStatus(t *testing.T) {
	storage := NewSaleStorage()
	service := NewSaleService(storage)

	// Crear algunas ventas de prueba
	sales := []Sale{
		{
			Id:        "sale-1",
			UserId:    "user-1",
			Amount:    100.0,
			Status:    Pending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
		{
			Id:        "sale-2",
			UserId:    "user-1",
			Amount:    200.0,
			Status:    Aproved,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
		{
			Id:        "sale-3",
			UserId:    "user-2",
			Amount:    300.0,
			Status:    Pending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
	}

	for _, s := range sales {
		storage.PutSale(&s)
	}

	t.Run("Obtener ventas por usuario y estado", func(t *testing.T) {
		result, err := service.GetByUserStatus("user-1", Pending)
		assert.NoError(t, err)
		assert.Len(t, *result, 1)
		assert.Equal(t, Pending, (*result)[0].Status)
	})

	t.Run("Obtener todas las ventas de un usuario", func(t *testing.T) {
		result, err := service.GetByUserStatus("user-1", "")
		assert.NoError(t, err)
		assert.Len(t, *result, 2)
	})

	t.Run("Obtener ventas con estado inválido", func(t *testing.T) {
		result, err := service.GetByUserStatus("user-1", "estado_invalido")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidStatus, err)
		assert.Nil(t, result)
	})
}
