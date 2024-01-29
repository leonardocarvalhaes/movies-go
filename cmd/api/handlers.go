package main

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) SaveMovie(w http.ResponseWriter, r *http.Request) {
	movie, err := app.fromRequestToMovie(w, r)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = func() error {
		if movie.ID > 0 {
			return app.DB.UpdateMovie(movie)
		}

		return app.DB.SaveMovie(movie)
	}()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := dtos.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Movie %d successfuly saved", movie.ID),
	}

	_ = app.writeJSON(w, http.StatusCreated, response)
}

func (app *application) ImportMovies(w http.ResponseWriter, r *http.Request) {
	var importParams dtos.ImportMovies

	err := app.readJSON(w, r, &importParams)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	otherWords := strings.Split(importParams.OtherWords, " ")

	var searchTerms []string

	for _, otherWord := range otherWords {
		searchTerms = append(searchTerms, importParams.PivotWords+" "+otherWord)
	}

	var movies []models.Movie
	var importedIDs []int

	isTheIDAlreadyImported := func(importedList []int, newID int) bool {
		for _, importedID := range importedIDs {
			if newID == importedID {
				return true
			}
		}

		return false
	}

	for _, searchTerm := range searchTerms {
		tmdbResponse, err := app.searchTMDB(searchTerm)

		if err != nil {
			log.Println(err)
		}

		for _, tmdbResult := range tmdbResponse.Results {
			ratingInt, _ := strconv.Atoi(importParams.Rating)
			voteCountInt, _ := strconv.Atoi(importParams.Votes)
			yearInt, _ := strconv.Atoi(importParams.Year)

			ratingIsntEnough := tmdbResult.VoteAverage < float32(ratingInt)
			voteCountIsntEnough := float32(tmdbResult.VoteCount) < float32(voteCountInt)

			releaseDate, err := time.Parse("2006-01-02", tmdbResult.ReleaseDate)

			if err != nil {
				fmt.Println(err)
				continue
			}

			if yearInt > releaseDate.Year() || ratingIsntEnough || voteCountIsntEnough || tmdbResult.PosterPath == "" {
				continue
			}

			if isTheIDAlreadyImported(importedIDs, tmdbResult.ID) {
				continue
			}

			moviesInDB, _ := app.DB.GetMovies(dtos.SearchColumn{
				Column:   "title",
				Operator: "=",
				Value:    tmdbResult.Title,
			})

			if len(moviesInDB) > 0 {
				continue
			}

			movie := models.Movie{
				Title:       tmdbResult.Title,
				Description: tmdbResult.Overview,
				ReleaseDate: releaseDate,
				MPAARating:  "L",
				Rating:      tmdbResult.VoteAverage,
				VoteCount:   tmdbResult.VoteCount,
				Image:       tmdbResult.PosterPath,
			}

			movies = append(movies, movie)
			importedIDs = append(importedIDs, tmdbResult.ID)
		}
	}

	for _, movie := range movies {
		_ = app.DB.SaveMovie(&movie)
	}

	response := dtos.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Successfuly imported %d movies", len(movies)),
	}

	_ = app.writeJSON(w, http.StatusCreated, response)
}

func (app *application) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.GetAllMovies()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.GetMovieByID(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)
}

func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteMovie(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = dtos.JSONResponse{
		Error:   false,
		Message: "Successfuly executed",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) GetMoviesCatalogue(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.GetAllMovies()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.DB.GetGenres()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, genres)
}

func (app *application) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.DB.GetUserByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	tokenPair, err := app.auth.GenerateTokenPair(user)

	if err != nil {
		app.errorJSON(w, err)
	}

	http.SetCookie(w, app.auth.GetRefreshCookie(tokenPair.RefreshToken))

	app.writeJSON(w, http.StatusOK, tokenPair)
}

func (app *application) RefreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(app.JWTSecret), nil
			})

			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			userID, err := strconv.Atoi(claims.Subject)

			if err != nil {
				app.errorJSON(w, errors.New("unkown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)

			if err != nil {
				app.errorJSON(w, errors.New("unkown user"), http.StatusUnauthorized)
				return
			}

			tokenPair, err := app.auth.GenerateTokenPair(user)

			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPair.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPair)

			return
		}
	}

	app.errorJSON(w, errors.New("refresh token absent"), http.StatusUnauthorized)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			cookie.MaxAge = -1
			cookie.Expires = time.Now().Add(-150 * time.Hour)
			http.SetCookie(w, cookie)
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
