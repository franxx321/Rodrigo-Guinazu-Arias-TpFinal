package api

import (
	"Rodrigo-Guinazu-Arias-TpFinal/models/Sales"
	"Rodrigo-Guinazu-Arias-TpFinal/models/users"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService *users.UserService
}

type SaleHandler struct {
	saleService *Sales.SaleService
}

func (h *UserHandler) HandleUserCreate(ctx *gin.Context) {
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

func (h *UserHandler) HandleUserRead(ctx *gin.Context) {
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

func (h *UserHandler) HandleUserUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var fields *users.UpdateFields
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//FIXME hay que modificar esto en base a la modificacion que hizo el joako en el userService
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
}*/

func (h *UserHandler) HandleUserDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.userService.Delete(id); err != nil {
		if errors.Is(err, users.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *SaleHandler) HandleSaleCreate(ctx *gin.Context) {
	var req struct {
		UserId string  `json:"user_id"`
		Amount float32 `json:"amount"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sale, err := h.saleService.Create(req.UserId, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, sale)
	}
	return
}

func (h *SaleHandler) HandleSaleUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.saleService.Update(id, req.Status)
	if err != nil {
		if errors.Is(err, Sales.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, Sales.ErrSaleNotPending) || errors.Is(err, Sales.ErrInvalidTransition) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, Sales.ErrInvalidStatus) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Sale updated"})
	}
	return
}

func (h *SaleHandler) HandleSaleGetByUserStatus(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	status := ctx.Query("status")

	sales, err := h.saleService.GetByUserStatus(userId, status)
	if err != nil {
		if errors.Is(err, Sales.ErrInvalidStatus) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var quantity int
	var aproved int
	var rejected int
	var pending int
	var totalAmount float32
	for _, sale := range *sales {
		quantity++
		if sale.Status == Sales.Aproved {
			aproved++
		} else if sale.Status == Sales.Rejected {
			rejected++
		} else if sale.Status == Sales.Pending {
			pending++
		}
		totalAmount += sale.Amount
	}
	var response struct {
		Metadata struct {
			Quantity    int     `json:"quantity"`
			Aproved     int     `json:"aproved"`
			Rejected    int     `json:"rejected"`
			Pending     int     `json:"pending"`
			TotalAmount float32 `json:"total_amount"`
		} `json:"metadata"`
		Sales []Sales.Sale `json:"results"`
	}
	ctx.JSON(http.StatusOK, response)
}
