package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"movie-rent/pkg/movie/mocks"
	"movie-rent/pkg/movie/model"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type MovieControllerTestSuite struct {
	suite.Suite
	context          *gin.Context
	recorder         *httptest.ResponseRecorder
	mockController   *gomock.Controller
	mockMovieService *mocks.MockMovieService
	testController   MovieController
}

func TestMovieControllerTestSuite(t *testing.T) {
	suite.Run(t, new(MovieControllerTestSuite))
}

func (suite *MovieControllerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockMovieService = mocks.NewMockMovieService(suite.mockController)
	suite.testController = NewMovieController(suite.mockMovieService)
}

func (suite *MovieControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *MovieControllerTestSuite) Test_AddMovie_ShouldReturnInternalServerErrorWhenServiceReturnError() {
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/movie", nil)
	suite.mockMovieService.EXPECT().AddMovie().Return(errors.New("error")).Times(1)

	suite.testController.AddMovie(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_AddMovie_ShouldStoreAllMoviesToDB() {
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/movie", nil)
	suite.mockMovieService.EXPECT().AddMovie().Return(nil).Times(1)

	suite.testController.AddMovie(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovies_ShouldReturnInternalServerErrorWhenServiceReturnError() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies", nil)
	suite.mockMovieService.EXPECT().GetMovies().Return([]model.Movie{}, errors.New("error")).Times(1)

	suite.testController.GetMovies(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovies_ShouldReturnAllMoviesSuccessfully() {
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
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies", nil)
	suite.mockMovieService.EXPECT().GetMovies().Return(movies, nil).Times(1)

	suite.testController.GetMovies(suite.context)

	expectedMovies := `[{"id":1,"title":"Hero","releaseYear":1990,"genre":"Action","description":"Action movie","imdbCode":"1234"}]`
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(expectedMovies, suite.recorder.Body.String())
}

func (suite *MovieControllerTestSuite) Test_GetFilteredMovies_ShouldReturnBadRequestWhenSearchTypeIsEmpty() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies/filter?searchType=&searchText=action", nil)

	suite.testController.GetFilteredMovies(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetFilteredMovies_ShouldReturnBadRequestWhenSearchTextIsEmpty() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies/filter?searchType=genre&searchText=", nil)

	suite.testController.GetFilteredMovies(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetFilteredMovies_ShouldReturnErrorServiceCallFailed() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies/filter?searchType=genre&searchText=action", nil)
	suite.mockMovieService.EXPECT().GetFilteredMovies("genre", "action").Return(nil, errors.New("error")).Times(1)

	suite.testController.GetFilteredMovies(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetFilteredMovies_ShouldReturnFilteredMovies() {
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
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies/filter?searchType=genre&searchText=action", nil)
	suite.mockMovieService.EXPECT().GetFilteredMovies("genre", "action").Return(movies, nil).Times(1)

	suite.testController.GetFilteredMovies(suite.context)

	expectedMovies := `[{"id":1,"title":"Hero","releaseYear":1990,"genre":"Action","description":"Action movie","imdbCode":"1234"}]`
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(expectedMovies, suite.recorder.Body.String())
}

func (suite *MovieControllerTestSuite) Test_GetMovieBy_ShouldReturnBadRequestWhenMovieIdNotPresentInPathParams() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))

	suite.testController.GetMovieBy(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieBy_ShouldReturnBadRequestWhenMovieIdNotANumber() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: "user",
		},
	}

	suite.testController.GetMovieBy(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieBy_ShouldReturnInternalServerErrorWhenServiceCallFailed() {
	movieId := 1001
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: strconv.Itoa(movieId),
		},
	}
	suite.mockMovieService.EXPECT().GetMovieBy(movieId).Return(model.Movie{}, errors.New("error")).Times(1)

	suite.testController.GetMovieBy(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieBy_ShouldSuccessfullyAddMovieToCart() {
	movieId := 1001
	expectedResponse := model.Movie{
		Id:          1,
		Title:       "Hero",
		Year:        1990,
		Genre:       "Action",
		Description: "Action movie",
		ImdbCode:    "1234",
	}
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: strconv.Itoa(movieId),
		},
	}
	suite.mockMovieService.EXPECT().GetMovieBy(movieId).Return(expectedResponse, nil).Times(1)

	suite.testController.GetMovieBy(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}
