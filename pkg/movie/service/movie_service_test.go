package service

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"movie-rent/pkg/movie/mocks"
	"movie-rent/pkg/movie/model"
	"testing"
)

type MovieServiceTestSuite struct {
	suite.Suite
	mockRapidClient *mocks.MockRapidClient
	mockController  *gomock.Controller
	mockRepository  *mocks.MockMovieRepository

	movieService MovieService
}

func TestMovieServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MovieServiceTestSuite))
}

func (suite *MovieServiceTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepository = mocks.NewMockMovieRepository(suite.mockController)
	suite.mockRapidClient = mocks.NewMockRapidClient(suite.mockController)

	suite.movieService = NewMovieService(suite.mockRepository, suite.mockRapidClient)
}

func (suite *MovieServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *MovieServiceTestSuite) Test_ReturnErrorWhenRapidClientFailed() {
	suite.mockRapidClient.EXPECT().FetchAllMovies().Return(nil, fmt.Errorf("error")).Times(1)

	err := suite.movieService.AddMovie()

	suite.NotNil(err)
}

func (suite *MovieServiceTestSuite) Test_ReturnErrorWhileStoringMovieToDB() {
	movies := []model.Movie{
		{
			Id:          1,
			Title:       "Hero",
			Year:        1990,
			Genre:       "Action",
			Description: "Action movie",
			ImdbCode:    "1234",
		},
	}
	suite.mockRapidClient.EXPECT().FetchAllMovies().Return(movies, nil).Times(1)

	suite.mockRepository.EXPECT().SaveAll(movies).Return(fmt.Errorf("error"))

	err := suite.movieService.AddMovie()

	suite.NotNil(err)
}

func (suite *MovieServiceTestSuite) Test_ShouldStoreMovieToDBSuccessfully() {
	movies := []model.Movie{
		{
			Id:          1,
			Title:       "Hero",
			Year:        1990,
			Genre:       "Action",
			Description: "Action movie",
			ImdbCode:    "1234",
		},
	}
	suite.mockRapidClient.EXPECT().FetchAllMovies().Return(movies, nil).Times(1)

	suite.mockRepository.EXPECT().SaveAll(movies).Return(nil)

	err := suite.movieService.AddMovie()

	suite.Nil(err)
}

func (suite *MovieServiceTestSuite) Test_ReturnErrorFetchToFailedMovies() {
	suite.mockRepository.EXPECT().GetMovies().Return([]model.Movie{}, fmt.Errorf("error"))

	movies, err := suite.movieService.GetMovies()

	suite.NotNil(err)
	suite.Equal(0, len(movies))
}

func (suite *MovieServiceTestSuite) Test_ShouldReturnMovies() {
	expectedMovies := []model.Movie{
		{
			Id:          1,
			Title:       "Hero",
			Year:        1990,
			Genre:       "Action",
			Description: "Action movie",
			ImdbCode:    "1234",
		},
	}
	suite.mockRepository.EXPECT().GetMovies().Return(expectedMovies, nil).Times(1)

	movies, err := suite.movieService.GetMovies()

	suite.Nil(err)
	suite.Equal(expectedMovies, movies)
}

func (suite *MovieServiceTestSuite) Test_GetFilteredMovies_ReturnErrorFetchToFailedMovies() {
	searchType := "searchType"
	searchText := "searchText"
	suite.mockRepository.EXPECT().FetchMoviesBySearchText(searchType, searchText).Return([]model.Movie{}, fmt.Errorf("error"))

	movies, err := suite.movieService.GetFilteredMovies(searchType, searchText)

	suite.NotNil(err)
	suite.Equal(0, len(movies))
}

func (suite *MovieServiceTestSuite) Test_GetFilteredMovies_ShouldReturnFilterMoviesByYear() {
	searchType := "year"
	searchText := "1990"
	expectedMovies := []model.Movie{
		{
			Id:          1,
			Title:       "Hero",
			Year:        1990,
			Genre:       "Action",
			Description: "Action movie",
			ImdbCode:    "1234",
		},
	}
	suite.mockRepository.EXPECT().FetchMoviesByYear(1990).Return(expectedMovies, nil).Times(1)

	movies, err := suite.movieService.GetFilteredMovies(searchType, searchText)

	suite.Nil(err)
	suite.Equal(expectedMovies, movies)
}

func (suite *MovieServiceTestSuite) Test_GetFilteredMovies_ShouldReturnFilterMoviesBySearchText() {
	searchType := "Genre"
	searchText := "Action"
	expectedMovies := []model.Movie{
		{
			Id:          1,
			Title:       "Hero",
			Year:        1990,
			Genre:       "Action",
			Description: "Action movie",
			ImdbCode:    "1234",
		},
	}
	suite.mockRepository.EXPECT().FetchMoviesBySearchText(searchType, searchText).Return(expectedMovies, nil).Times(1)

	movies, err := suite.movieService.GetFilteredMovies(searchType, searchText)

	suite.Nil(err)
	suite.Equal(expectedMovies, movies)
}

func (suite *MovieServiceTestSuite) Test_AddMovieToCart_ShouldReturnErrorAddMovieToCartFailed() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	suite.mockRepository.EXPECT().AddMovieToCart(request).Return(fmt.Errorf("error"))

	err := suite.movieService.AddMovieToCart(request)

	suite.NotNil(err)
}

func (suite *MovieServiceTestSuite) Test_AddMovieToCart_ShouldSuccessfullyAddMovieToCart() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	suite.mockRepository.EXPECT().AddMovieToCart(request).Return(nil).Times(1)

	err := suite.movieService.AddMovieToCart(request)

	suite.Nil(err)
}
