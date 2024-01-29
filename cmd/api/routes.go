package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/movies", app.GetMovies)
	mux.Get("/movies/{id}", app.GetMovie)
	mux.Get("/genres", app.GetGenres)
	mux.Post("/authenticate", app.Authenticate)
	mux.Get("/refresh", app.RefreshToken)
	mux.Get("/logout", app.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/catalogue", app.GetMoviesCatalogue)
		mux.Put("/movies/create", app.SaveMovie)
		mux.Post("/movies/import", app.ImportMovies)
		mux.Patch("/movies/{id}", app.SaveMovie)
		mux.Delete("/movies/{id}", app.DeleteMovie)
	})

	return mux
}
