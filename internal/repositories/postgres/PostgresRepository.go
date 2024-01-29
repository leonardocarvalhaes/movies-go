package postgres

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"context"
	"database/sql"
	"strings"
	"time"
)

type PostgresRepository struct {
	DB *sql.DB
}

const connectionTimeout = time.Second * 3

func (r *PostgresRepository) GetConnection() *sql.DB {
	return r.DB
}

func (r *PostgresRepository) SaveMovie(movie *models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `
		insert into movies
			(title, release_date, runtime,
			mpaa_rating, rating, vote_count,
			description, image, created_at, updated_at)
		values
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id
	`

	row := r.DB.QueryRowContext(ctx, query,
		movie.Title,
		movie.ReleaseDate,
		movie.Duration,
		movie.MPAARating,
		movie.Rating,
		movie.VoteCount,
		movie.Description,
		movie.Image,
		time.Now(),
		time.Now(),
	)

	err := row.Scan(
		&movie.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) UpdateMovie(movie *models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `
		update movies
			set title = $1, release_date = $2,
			runtime = $3, mpaa_rating = $4, description = $5,
			image = $6, updated_at=$7
		where
			id = $8
	`

	r.DB.ExecContext(ctx, query,
		movie.Title,
		movie.ReleaseDate,
		movie.Duration,
		movie.Rating,
		movie.Description,
		movie.Image,
		time.Now(),
		movie.ID,
	)

	return nil
}

func (r *PostgresRepository) GetAllMovies() ([]*models.Movie, error) {
	return r.GetMovies()
}

func (r *PostgresRepository) GetMovies(searchColumns ...dtos.SearchColumn) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	var whereConditions []string

	for _, searchColumn := range searchColumns {
		whereConditions = append(whereConditions, searchColumn.Column+searchColumn.Operator+"'"+searchColumn.Value+"'")
	}

	where := strings.Join(whereConditions, " and ")

	queryStart := `
		select
			id, title, release_date, runtime,
			mpaa_rating, coalesce(rating, 0.0), coalesce(vote_count, 0), description,
			coalesce(image, ''), created_at, updated_at
		from
			movies
			`
	queryEnd := `order by title`

	query := queryStart + " " + where + " " + queryEnd

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
			&movie.MPAARating,
			&movie.Rating,
			&movie.VoteCount,
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

func (r *PostgresRepository) GetMovieByID(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `
		select
			id, title, release_date, runtime,
			mpaa_rating, coalesce(rating, 0.0), coalesce(vote_count, 0), description,
			coalesce(image, ''), created_at, updated_at
		from
			movies
		where
			id = $1
		order by
			title
	`

	row := r.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Duration,
		&movie.MPAARating,
		&movie.Rating,
		&movie.VoteCount,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *PostgresRepository) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	query := `delete from movies where id = $1`

	_, err := r.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetGenres(searchColumns ...dtos.SearchColumn) ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	var whereConditions []string

	for _, searchColumn := range searchColumns {
		whereConditions = append(whereConditions, searchColumn.Column+searchColumn.Operator+"'"+searchColumn.Value+"'")
	}

	where := strings.Join(whereConditions, " and ")

	queryStart := `
		select
			id, genre, created_at, updated_at
		from
			genres
			`
	queryEnd := `order by genre`

	query := queryStart + " " + where + " " + queryEnd

	rows, err := r.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var genre models.Genre

		err := rows.Scan(
			&genre.ID,
			&genre.Name,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		genres = append(genres, &genre)
	}

	return genres, nil
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
