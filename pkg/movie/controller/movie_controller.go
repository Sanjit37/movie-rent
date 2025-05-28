package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
