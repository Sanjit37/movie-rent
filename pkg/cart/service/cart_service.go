package service

import (
	"fmt"
	"movie-rent/pkg/cart/model"
	"movie-rent/pkg/cart/repository"
)

// go:generate mockgen -source=pkg/cart/service/cart_service.go -destination=pkg/cart/mocks/cart_service_mock.go -package=mocks

type CartService interface {
	AddToCart(request model.CartRequest) (int, error)
	GetCartItems(userId int) ([]model.CartResponse, error)
}

type cartService struct {
	repository repository.CartRepository
}

func NewCartService(repository repository.CartRepository) CartService {
	return cartService{repository: repository}
}

func (m cartService) AddToCart(request model.CartRequest) (int, error) {
	id, err := m.repository.AddToCart(request)
	if err != nil {
		fmt.Println("failed to add to cart: %w", err.Error())
		return 0, err
	}
	return id, nil
}

func (m cartService) GetCartItems(userId int) ([]model.CartResponse, error) {
	res, err := m.repository.GetCartItems(userId)
	if err != nil {
		fmt.Println("failed to add to cart: %w", err.Error())
		return []model.CartResponse{}, err
	}
	return res, nil
}
