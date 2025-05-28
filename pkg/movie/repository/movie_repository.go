package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"movie-rent/pkg/movie/model"
)

const (
	InsertMovieSQL = `INSERT INTO movies(id, title, description, gener, release_year, imdb_code) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	SelectMovies   = `SELECT id, title, release_year, gener, description, imdb_code FROM movies`
)

type MovieRepository interface {
	Save(movie model.Movie) error
	SaveAll(movies []model.Movie) error
	GetMovies() ([]model.Movie, error)
}

type movieRepo struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) MovieRepository {
	return &movieRepo{db: db}
}

func (m movieRepo) Save(movie model.Movie) error {
	id := movie.Id
	err := m.db.QueryRow(InsertMovieSQL, movie.Id, movie.Title, movie.Description, movie.Gener, movie.Year, movie.ImdbCode).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}
	fmt.Println("movie inserted. Id:", id)
	return nil
}

func (m movieRepo) SaveAll(movies []model.Movie) error {
	for _, movie := range movies {
		err := m.Save(movie)
		if err != nil {
			fmt.Println("failed to save movie: %w", err.Error())
		}
	}
	return nil
}

func (m movieRepo) GetMovies() ([]model.Movie, error) {
	rows, err := m.db.Query(SelectMovies)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Gener, &movie.Description, &movie.ImdbCode)
		if err != nil {
			log.Println("Error scanning row:", err)
		}
		movies = append(movies, movie)
	}

	fmt.Println("movies fetched. Total movies:", len(movies))
	return movies, nil
}
