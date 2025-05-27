package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *UserService) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	storage := NewUserStorage()
	service := NewUserService(storage)

	router.POST("/users", func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`     // Agregamos validación
			Address  string `json:"address" binding:"required"`  // Agregamos validación
			NickName string `json:"nickname" binding:"required"` // Agregamos validación
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u := &User{
			Name:     req.Name,
			Address:  req.Address,
			NickName: req.NickName,
		}

		if err := service.Create(u); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, u)
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		user, err := service.Get(id)
		if err != nil {
			if err == ErrNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	return router, service
}

func TestUserService_Create(t *testing.T) {
	router, _ := setupRouter()

	t.Run("Crear usuario exitosamente", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/users", strings.NewReader(`{
			"name": "John Doe",
			"address": "123 Main St",
			"nickname": "johndoe"
		}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response User
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, "John Doe", response.Name)
		assert.Equal(t, "123 Main St", response.Address)
		assert.Equal(t, "johndoe", response.NickName)
		assert.Equal(t, 1, response.Version)
		assert.False(t, response.CreatedAt.IsZero())
		assert.False(t, response.UpdatedAt.IsZero())
	})

	t.Run("Crear usuario con datos inválidos", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/users", strings.NewReader(`{
			"name": ""
		}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserService_Get(t *testing.T) {
	router, service := setupRouter()

	t.Run("Obtener usuario existente", func(t *testing.T) {
		// Crear un usuario primero
		user := &User{
			Name:     "Test User",
			Address:  "Test Address",
			NickName: "testuser",
		}
		err := service.Create(user)
		assert.NoError(t, err)

		// Intentar obtener el usuario
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/"+user.ID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response User
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
	})

	t.Run("Obtener usuario inexistente", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/users/nonexistent", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserService_Update(t *testing.T) {
	storage := NewUserStorage()
	service := NewUserService(storage)

	t.Run("Actualizar usuario existente", func(t *testing.T) {
		// Crear un usuario
		user := &User{
			Name:     "Original Name",
			Address:  "Original Address",
			NickName: "original",
		}
		err := service.Create(user)
		assert.NoError(t, err)

		// Esperar un momento para asegurar que UpdatedAt sea diferente
		time.Sleep(time.Millisecond * 100)

		// Preparar campos para actualizar
		newName := "Updated Name"
		updateFields := &UpdateFields{
			Name: &newName,
		}

		// Actualizar usuario
		updatedUser, err := service.Update(user.ID, updateFields)
		assert.NoError(t, err)
		assert.Equal(t, newName, updatedUser.Name)
		assert.Equal(t, user.Address, updatedUser.Address)
		assert.Equal(t, 2, updatedUser.Version)
		assert.True(t, updatedUser.UpdatedAt.After(user.CreatedAt))
	})

	t.Run("Actualizar usuario inexistente", func(t *testing.T) {
		newName := "Updated Name"
		updateFields := &UpdateFields{
			Name: &newName,
		}

		_, err := service.Update("nonexistent", updateFields)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestUserService_Delete(t *testing.T) {
	storage := NewUserStorage()
	service := NewUserService(storage)

	t.Run("Eliminar usuario existente", func(t *testing.T) {
		// Crear un usuario
		user := &User{
			Name:     "To Delete",
			Address:  "Delete Address",
			NickName: "delete",
		}
		err := service.Create(user)
		assert.NoError(t, err)

		// Eliminar usuario
		err = service.Delete(user.ID)
		assert.NoError(t, err)

		// Verificar que el usuario fue eliminado
		_, err = service.Get(user.ID)
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("Eliminar usuario inexistente", func(t *testing.T) {
		err := service.Delete("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}
