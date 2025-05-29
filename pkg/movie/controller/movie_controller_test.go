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

func (suite *MovieControllerTestSuite) Test_AddMovieToCart_ShouldReturnBadRequestWhenRequiredFieldIsEmpty() {
	invalidRequestBody := `{"userId":1001,"movieId":4563,"movieName":"Hero"}`
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/addMovieToCart", strings.NewReader(invalidRequestBody))

	suite.testController.AddMovieToCart(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_AddMovieToCart_ShouldReturnInternalServerErrorWhenServiceCallFailed() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	requestBody := `{"userId":1001,"movieId":4563,"movieName":"Hero","releaseYear":1990}`
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/addMovieToCart", strings.NewReader(requestBody))
	suite.mockMovieService.EXPECT().AddMovieToCart(request).Return(0, errors.New("error")).Times(1)

	suite.testController.AddMovieToCart(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_AddMovieToCart_ShouldSuccessfullyAddMovieToCart() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	requestBody := `{"userId":1001,"movieId":4563,"movieName":"Hero","releaseYear":1990}`
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/addMovieToCart", strings.NewReader(requestBody))
	suite.mockMovieService.EXPECT().AddMovieToCart(request).Return(1, nil).Times(1)

	suite.testController.AddMovieToCart(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(1, suite.recorder.Body.Len())
}

func (suite *MovieControllerTestSuite) Test_GetCartList_ShouldReturnBadRequestWhenUserIdNotPresentInPathParams() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cartList/1001", strings.NewReader(""))

	suite.testController.GetCartList(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetCartList_ShouldReturnBadRequestWhenUserIdNotANumber() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cartList/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "userId",
			Value: "user",
		},
	}

	suite.testController.GetCartList(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetCartList_ShouldReturnInternalServerErrorWhenServiceCallFailed() {
	userId := 1001
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cartList/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "userId",
			Value: strconv.Itoa(userId),
		},
	}
	suite.mockMovieService.EXPECT().GetCartList(userId).Return(nil, errors.New("error")).Times(1)

	suite.testController.GetCartList(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetCartList_ShouldSuccessfullyAddMovieToCart() {
	userId := 1001
	expectedResponse := []model.CartResponse{{
		Id:          1,
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}}
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cartList", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "userId",
			Value: strconv.Itoa(userId),
		},
	}
	suite.mockMovieService.EXPECT().GetCartList(userId).Return(expectedResponse, nil).Times(1)

	suite.testController.GetCartList(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieDetailsBy_ShouldReturnBadRequestWhenMovieIdNotPresentInPathParams() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))

	suite.testController.GetMovieDetailsBy(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieDetailsBy_ShouldReturnBadRequestWhenMovieIdNotANumber() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: "user",
		},
	}

	suite.testController.GetMovieDetailsBy(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieDetailsBy_ShouldReturnInternalServerErrorWhenServiceCallFailed() {
	movieId := 1001
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/movie/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: strconv.Itoa(movieId),
		},
	}
	suite.mockMovieService.EXPECT().GetMovieDetailsBy(movieId).Return(model.Movie{}, errors.New("error")).Times(1)

	suite.testController.GetMovieDetailsBy(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetMovieDetailsBy_ShouldSuccessfullyAddMovieToCart() {
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
	suite.mockMovieService.EXPECT().GetMovieDetailsBy(movieId).Return(expectedResponse, nil).Times(1)

	suite.testController.GetMovieDetailsBy(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}
