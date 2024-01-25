package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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

func (app *application) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.GetAllMovies()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) GetMoviesCatalogue(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.GetAllMovies()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
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
