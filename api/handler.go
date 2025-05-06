package api

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	userService *users.UserService
}

func (h *handler) handleCreate(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		NickName string `json:"nickname"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &users.User{
		Name:     req.Name,
		Address:  req.Address,
		NickName: req.NickName,
	}

	if err := h.userService.Create(u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

func (h *handler) handleUserRead(ctx *gin.Context) {
	id := ctx.Param("id")

	u, err := h.userService.Get(id)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (h *handler) handleUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var fields *users.UpdateFields
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.userService.Update(id, fields)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, u)
}

func (h *handler) handleDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.userService.Delete(id); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
