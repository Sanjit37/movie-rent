package model

type CartRequest struct {
	UserId      int    `json:"userId"  binding:"required"`
	MovieId     int    `json:"movieId"  binding:"required"`
	MovieName   string `json:"movieName"  binding:"required"`
	ReleaseYear int    `json:"releaseYear"  binding:"required"`
}
type CartResponse struct {
	Id          int    `json:"id"`
	UserId      int    `json:"userId"  binding:"required"`
	MovieId     int    `json:"movieId"  binding:"required"`
	MovieName   string `json:"movieName"  binding:"required"`
	ReleaseYear int    `json:"releaseYear"  binding:"required"`
}
