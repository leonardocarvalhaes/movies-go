package models

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Duration    int       `json:"duration"`
	MPAARating  string    `json:"mpaa_rating"`
	Rating      float32   `json:"rating"`
	VoteCount   int       `json:"vote_count"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
