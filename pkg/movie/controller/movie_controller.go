package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"movie-rent/pkg/movie/service"
	"net/http"
	"strconv"
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
		ctx.JSON(http.StatusBadRequest, "searchType or searchText is empty")
		return
	}
	
	movies, err := m.service.GetFilteredMovies(searchType, searchText)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, movies)
}

func (m *MovieController) GetMovieBy(ctx *gin.Context) {
	fmt.Println("Fetching movie by id")
	id := ctx.Param("id")
	movieId, err := strconv.Atoi(id)
	if id == "" || err != nil {
		ctx.JSON(http.StatusBadRequest, "Invalid id")
		return
	}

	movie, serviceErr := m.service.GetMovieBy(movieId)
	if serviceErr != nil {
		ctx.JSON(http.StatusInternalServerError, serviceErr)
		return
	}
	ctx.JSON(http.StatusOK, movie)
}
