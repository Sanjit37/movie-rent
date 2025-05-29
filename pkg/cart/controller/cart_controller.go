package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"movie-rent/pkg/cart/model"
	"movie-rent/pkg/cart/service"
	"net/http"
	"strconv"
)

type CartController struct {
	service service.CartService
}

func NewCartController(service service.CartService) CartController {
	return CartController{service: service}
}

func (m *CartController) AddToCart(ctx *gin.Context) {
	var cart model.CartRequest
	bindErr := ctx.ShouldBindJSON(&cart)
	if bindErr != nil {
		fmt.Println("Invalid request body", bindErr.Error())
		ctx.JSON(http.StatusBadRequest, bindErr.Error())
		return
	}
	id, err := m.service.AddToCart(cart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func (m *CartController) GetCartItems(ctx *gin.Context) {
	var userId int
	var err error
	id := ctx.Param("userId")
	userId, err = strconv.Atoi(id)
	if id == "" || err != nil {
		ctx.JSON(http.StatusBadRequest, "userId is empty")
		return
	}
	res, err := m.service.GetCartItems(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}
