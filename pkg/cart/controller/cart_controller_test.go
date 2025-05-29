package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"movie-rent/pkg/cart/mocks"
	"movie-rent/pkg/cart/model"
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
	mockMovieService *mocks.MockCartService
	testController   CartController
}

func TestMovieControllerTestSuite(t *testing.T) {
	suite.Run(t, new(MovieControllerTestSuite))
}

func (suite *MovieControllerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockMovieService = mocks.NewMockCartService(suite.mockController)
	suite.testController = NewCartController(suite.mockMovieService)
}

func (suite *MovieControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *MovieControllerTestSuite) Test_AddToCart_ShouldReturnBadRequestWhenRequiredFieldIsEmpty() {
	invalidRequestBody := `{"userId":1001,"movieId":4563,"movieName":"Hero"}`
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/cart/add", strings.NewReader(invalidRequestBody))

	suite.testController.AddToCart(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_AddToCart_ShouldReturnInternalServerErrorWhenServiceCallFailed() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	requestBody := `{"userId":1001,"movieId":4563,"movieName":"Hero","releaseYear":1990}`
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/cart/add", strings.NewReader(requestBody))
	suite.mockMovieService.EXPECT().AddToCart(request).Return(0, errors.New("error")).Times(1)

	suite.testController.AddToCart(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_AddToCart_ShouldSuccessfullyAddToCart() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	requestBody := `{"userId":1001,"movieId":4563,"movieName":"Hero","releaseYear":1990}`
	suite.context.Request = httptest.NewRequest(http.MethodPost, "/cart/add", strings.NewReader(requestBody))
	suite.mockMovieService.EXPECT().AddToCart(request).Return(1, nil).Times(1)

	suite.testController.AddToCart(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(1, suite.recorder.Body.Len())
}

func (suite *MovieControllerTestSuite) Test_GetCartItems_ShouldReturnBadRequestWhenUserIdNotPresentInPathParams() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cart/items/1001", strings.NewReader(""))

	suite.testController.GetCartItems(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetCartItems_ShouldReturnBadRequestWhenUserIdNotANumber() {
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cart/items/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "userId",
			Value: "user",
		},
	}

	suite.testController.GetCartItems(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetCartItems_ShouldReturnInternalServerErrorWhenServiceCallFailed() {
	userId := 1001
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cart/items/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "userId",
			Value: strconv.Itoa(userId),
		},
	}
	suite.mockMovieService.EXPECT().GetCartItems(userId).Return(nil, errors.New("error")).Times(1)

	suite.testController.GetCartItems(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *MovieControllerTestSuite) Test_GetCartItems_ShouldSuccessfullyAddToCart() {
	userId := 1001
	expectedResponse := []model.CartResponse{{
		Id:          1,
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}}
	suite.context.Request = httptest.NewRequest(http.MethodGet, "/cart/items/1001", strings.NewReader(""))
	suite.context.Params = gin.Params{
		gin.Param{
			Key:   "userId",
			Value: strconv.Itoa(userId),
		},
	}
	suite.mockMovieService.EXPECT().GetCartItems(userId).Return(expectedResponse, nil).Times(1)

	suite.testController.GetCartItems(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}
