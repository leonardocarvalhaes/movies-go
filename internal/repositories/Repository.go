package repositories

import (
	"backend/internal/models"
	"database/sql"
)

type Repository interface {
	GetConnection() *sql.DB
	GetAllMovies() ([]*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}
