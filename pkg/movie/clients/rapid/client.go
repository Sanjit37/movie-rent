package rapid

import (
	"encoding/json"
	"fmt"
	"movie-rent/constants"
	"movie-rent/pkg/movie/model"
	"net/http"
)

type RapidClient interface {
	FetchAllMovies() ([]model.Movie, error)
}

type rapidClient struct {
	http *http.Client
}

func NewRapidClient(http *http.Client) RapidClient {
	return rapidClient{http: http}
}

func (r rapidClient) FetchAllMovies() ([]model.Movie, error) {
	url := constants.RapidBaseURL + constants.RapidPathURL
	resp, err := r.http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch all movies")
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch movies:", resp.Status)
		return nil, fmt.Errorf("failed to fetch movies: %s", resp.Status)
	}

	var movies []model.Movie
	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		fmt.Println("Failed to parse movies:", err)
		return nil, err
	}
	return movies, nil
}
