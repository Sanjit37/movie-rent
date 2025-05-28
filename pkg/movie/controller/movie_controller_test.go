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

func (suite *MovieControllerTestSuite) Test_AddMovieToDBFailed() {
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/movie", nil)
	suite.mockMovieService.EXPECT().AddMovie().Return(errors.New("error")).Times(1)

	suite.testController.AddMovie(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_AddMovieToDBSuccessfully() {
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/movie", nil)
	suite.mockMovieService.EXPECT().AddMovie().Return(nil).Times(1)

	suite.testController.AddMovie(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieToDBFailed() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movies", nil)
	suite.mockMovieService.EXPECT().GetMovies().Return([]model.Movie{}, errors.New("error")).Times(1)

	suite.testController.GetMovies(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieToDBSuccessfully() {
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
