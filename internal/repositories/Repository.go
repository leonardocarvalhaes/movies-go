package repositories

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"database/sql"
)

type Repository interface {
	GetConnection() *sql.DB
	SaveMovie(movie *models.Movie) error
	UpdateMovie(movie *models.Movie) error
	GetMovies(searchColumns ...dtos.SearchColumn) ([]*models.Movie, error)
	GetAllMovies() ([]*models.Movie, error)
	GetMovieByID(id int) (*models.Movie, error)
	DeleteMovie(id int) error
	GetGenres(searchColumns ...dtos.SearchColumn) ([]*models.Genre, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}
