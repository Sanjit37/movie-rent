package service

import (
	"fmt"
	"movie-rent/pkg/movie/clients/rapid"
	"movie-rent/pkg/movie/model"
	"movie-rent/pkg/movie/repository"
)

// go:generate mockgen -source=pkg/movie/service/movie_service.go -destination=pkg/movie/mocks/movie_service_mock.go -package=mocks

type MovieService interface {
	AddMovie() error
	GetMovies() ([]model.Movie, error)
}

type movieService struct {
	repository repository.MovieRepository
	client     rapid.RapidClient
}

func NewMovieService(repository repository.MovieRepository, client rapid.RapidClient) MovieService {
	return movieService{repository: repository, client: client}
}

func (m movieService) AddMovie() error {
	movies, err := m.client.FetchAllMovies()
	if err != nil {
		return fmt.Errorf("failed to fetch movies from rapid api: %w", err)
	}

	err = m.repository.SaveAll(movies)
	if err != nil {
		return fmt.Errorf("failed to save movies to rapid api: %s", err.Error())
	}

	fmt.Println("movie inserted successfully")
	return nil
}

func (m movieService) GetMovies() ([]model.Movie, error) {
	movies, err := m.repository.GetMovies()
	if err != nil {
		fmt.Println("failed to find movies: %w", err.Error())
		return []model.Movie{}, err
	}

	return movies, err
}
