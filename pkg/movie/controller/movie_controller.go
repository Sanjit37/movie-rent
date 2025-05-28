package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"movie-rent/pkg/movie/model"
	"movie-rent/pkg/movie/service"
	"net/http"
)

type MovieController struct {
	service service.MovieService
}

func NewMovieController(service service.MovieService) MovieController {
	return MovieController{service: service}
}

func (m *MovieController) AddMovie(ctx *gin.Context) {
	err := m.service.AddMovie()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (m *MovieController) GetMovies(ctx *gin.Context) {
	fmt.Println("Fetching all movies")
	movies, err := m.service.GetMovies()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

func (m *MovieController) GetFilteredMovies(ctx *gin.Context) {
	fmt.Println("Fetching filtered movies")
	searchType := ctx.Query("searchType")
	searchText := ctx.Query("searchText")
	if searchType == "" || searchText == "" {
		fmt.Println("Fetching filtered movies --1")
		ctx.JSON(http.StatusBadRequest, "searchType or searchText is empty")
		return
	}
	fmt.Println("Fetching filtered movies --2")
	movies, err := m.service.GetFilteredMovies(searchType, searchText)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, movies)
}

func (m *MovieController) AddMovieToCart(ctx *gin.Context) {
	var cart model.CartRequest
	err := ctx.ShouldBindJSON(&cart)
	if err != nil {
		fmt.Println("Invalid request body", err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = m.service.AddMovieToCart(cart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
