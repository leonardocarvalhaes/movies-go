package postgres

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostgresRepository struct {
	DB *sql.DB
}

const connectionTimeout = time.Second * 3

func (r *PostgresRepository) GetConnection() *sql.DB {
	return r.DB
}

func (r *PostgresRepository) GetAllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `
		select
			id, title, release_date, runtime,
			mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		from
			movies
		order by
			title
	`

	rows, err := r.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Duration,
			&movie.Rating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}

func (r *PostgresRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `
		select
			id, first_name, last_name, email,
			password, created_at, updated_at
		from
			users
		where
			email = $1
	`

	row := r.DB.QueryRowContext(ctx, query, email)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `
		select
			id, first_name, last_name, email,
			password, created_at, updated_at
		from
			users
		where
			id = $1
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
