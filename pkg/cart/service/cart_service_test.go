package service

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"movie-rent/pkg/cart/mocks"
	"movie-rent/pkg/cart/model"
	"testing"
)

type CartServiceTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
	mockRepository *mocks.MockCartRepository

	cartService CartService
}

func TestCartServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CartServiceTestSuite))
}

func (suite *CartServiceTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepository = mocks.NewMockCartRepository(suite.mockController)

	suite.cartService = NewCartService(suite.mockRepository)
}

func (suite *CartServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *CartServiceTestSuite) Test_AddToCart_ShouldReturnErrorAddToCartFailed() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	suite.mockRepository.EXPECT().AddToCart(request).Return(0, fmt.Errorf("error"))

	id, err := suite.cartService.AddToCart(request)

	suite.NotNil(err)
	suite.Equal(0, id)
}

func (suite *CartServiceTestSuite) Test_AddToCart_ShouldSuccessfullyAddToCart() {
	request := model.CartRequest{
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}
	suite.mockRepository.EXPECT().AddToCart(request).Return(1, nil).Times(1)

	id, err := suite.cartService.AddToCart(request)

	suite.Nil(err)
	suite.Equal(1, id)
}

func (suite *CartServiceTestSuite) Test_GetCartItems_ShouldReturnErrorWhenGetCartItemsFailed() {
	userId := 1001
	suite.mockRepository.EXPECT().GetCartItems(userId).Return(nil, fmt.Errorf("error")).Times(1)

	_, err := suite.cartService.GetCartItems(userId)

	suite.NotNil(err)
}

func (suite *CartServiceTestSuite) Test_GetCartItems_ShouldSuccessfullyFetchCartList() {
	userId := 1001
	response := []model.CartResponse{{
		Id:          1,
		UserId:      1001,
		MovieId:     4563,
		MovieName:   "Hero",
		ReleaseYear: 1990,
	}}
	suite.mockRepository.EXPECT().GetCartItems(userId).Return(response, nil).Times(1)

	actualResponse, err := suite.cartService.GetCartItems(userId)

	suite.Nil(err)
	suite.Equal(response, actualResponse)
}
