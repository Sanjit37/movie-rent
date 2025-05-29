package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"movie-rent/db"
	controller2 "movie-rent/pkg/cart/controller"
	repository2 "movie-rent/pkg/cart/repository"
	service2 "movie-rent/pkg/cart/service"
	"movie-rent/pkg/movie/clients/rapid"
	"movie-rent/pkg/movie/controller"
	"movie-rent/pkg/movie/repository"
	"movie-rent/pkg/movie/service"
	"net/http"
)

func main() {
	route := gin.Default()
	httpClient := &http.Client{}
	database := db.NewDatabase().Instance()

	movieRepository := repository.NewMovieRepository(database)
	rapidClient := rapid.NewRapidClient(httpClient)
	movieService := service.NewMovieService(movieRepository, rapidClient)
	movieController := controller.NewMovieController(movieService)

	cartRepository := repository2.NewCartRepository(database)
	cartService := service2.NewCartService(cartRepository)
	cartController := controller2.NewCartController(cartService)

	route.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Greetings": "Hello world"})
	})

	route.GET("/movies", movieController.GetMovies)
	route.GET("/movie/:id", movieController.GetMovieBy)
	route.GET("/movies/filter", movieController.GetFilteredMovies)
	route.POST("/movie", movieController.AddMovie)

	route.POST("/cart/add", cartController.AddToCart)
	route.GET("/cart/items/:userId", cartController.GetCartItems)

	route.Run(":8080")

	defer database.Close()
}
