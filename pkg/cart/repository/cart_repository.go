package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"movie-rent/pkg/cart/model"
)

const (
	InsertCartDetailsSQL = `INSERT INTO movie_carts(user_id, movie_id, movie_name, release_year) VALUES ($1, $2, $3, $4) RETURNING id`
	SelectCartListSQL    = `SELECT * FROM movie_carts where user_id = $1`
)

type CartRepository interface {
	AddToCart(cart model.CartRequest) (int, error)
	GetCartItems(userId int) ([]model.CartResponse, error)
}

type cartRepo struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &cartRepo{db: db}
}

func (m cartRepo) AddToCart(cart model.CartRequest) (int, error) {
	id := cart.MovieId
	err := m.db.QueryRow(InsertCartDetailsSQL, cart.UserId, cart.MovieId, cart.MovieName, cart.ReleaseYear).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert cart details: %w", err)
	}
	fmt.Println("Successfully inserted cart details. Id:", id)
	return id, nil
}

func (m cartRepo) GetCartItems(userId int) ([]model.CartResponse, error) {
	rows, err := m.db.Query(SelectCartListSQL, userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cartList []model.CartResponse
	for rows.Next() {
		var c model.CartResponse
		err := rows.Scan(&c.Id, &c.UserId, &c.MovieId, &c.MovieName, &c.ReleaseYear)
		if err != nil {
			log.Println("Error scanning row:", err)
		}
		cartList = append(cartList, c)
	}

	fmt.Println("Successfully fetched added movies", len(cartList))
	return cartList, nil
}
