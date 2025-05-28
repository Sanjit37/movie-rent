package model

type Movies struct {
	Movies Movie `json:"movies"`
}

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Year        int    `json:"releaseYear"`
	Gener       string `json:"gener"`
	Description string `json:"description"`
	ImdbCode    string `json:"imdbCode"`
}
