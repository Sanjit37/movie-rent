package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"movie-rent/pkg/movie/model"
)

const (
	InsertMovieSQL       = `INSERT INTO movies(id, title, description, genre, release_year, imdb_code) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	SelectMovies         = `SELECT id, title, release_year, genre, description, imdb_code FROM movies`
	SelectMoviesByYear   = `SELECT id, title, release_year, genre, description, imdb_code FROM movies where release_year = $1`
	InsertCartDetailsSQL = `INSERT INTO movies(user_id, movie_id, movie_name, release_year) VALUES ($1, $2, $3, $4) RETURNING id`
)

type MovieRepository interface {
	Save(movie model.Movie) error
	SaveAll(movies []model.Movie) error
	GetMovies() ([]model.Movie, error)
	FetchMoviesByYear(year int) ([]model.Movie, error)
	FetchMoviesBySearchText(searchType string, searchText string) ([]model.Movie, error)
	AddMovieToCart(cart model.CartRequest) error
}

type movieRepo struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) MovieRepository {
	return &movieRepo{db: db}
}

func (m movieRepo) Save(movie model.Movie) error {
	id := movie.Id
	err := m.db.QueryRow(InsertMovieSQL, movie.Id, movie.Title, movie.Description, movie.Genre, movie.Year, movie.ImdbCode).Scan(&id)

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
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Genre, &movie.Description, &movie.ImdbCode)
		if err != nil {
			log.Println("Error scanning row:", err)
		}
		movies = append(movies, movie)
	}

	fmt.Println("movies fetched. Total movies:", len(movies))
	return movies, nil
}

func (m movieRepo) FetchMoviesByYear(year int) ([]model.Movie, error) {
	rows, err := m.db.Query(SelectMoviesByYear, year)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Genre, &movie.Description, &movie.ImdbCode)
		if err != nil {
			log.Println("Error scanning row:", err)
		}
		movies = append(movies, movie)
	}

	fmt.Println("movies fetched. Total movies:", len(movies))
	return movies, nil
}

func (m movieRepo) FetchMoviesBySearchText(searchType string, searchText string) ([]model.Movie, error) {
	query := fmt.Sprintf("SELECT id, title, release_year, genre, description, imdb_code FROM movies WHERE %s ILIKE $1", searchType)
	rows, err := m.db.Query(query, "%"+searchText+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Genre, &movie.Description, &movie.ImdbCode)
		if err != nil {
			log.Println("Error scanning row:", err)
		}
		movies = append(movies, movie)
	}

	fmt.Println("movies fetched. Total movies:", len(movies))
	return movies, nil
}

func (m movieRepo) AddMovieToCart(cart model.CartRequest) error {
	id := cart.MovieId
	err := m.db.QueryRow(InsertCartDetailsSQL, cart.UserId, cart.MovieId, cart.MovieName, cart.ReleaseYear).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to insert cart details: %w", err)
	}
	fmt.Println("Successfully inserted cart details. Id:", id)
	return nil
}
